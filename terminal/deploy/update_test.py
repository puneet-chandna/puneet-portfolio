#!/usr/bin/env python3
"""Exercise release activation, rollback, and locking without touching Docker."""
import fcntl
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

SCRIPT = Path(__file__).with_name("update.sh")
OLD, NEW = "a" * 40 + "-1-1", "b" * 40 + "-2-1"

for previous, fail, locked in [(False, False, False), (True, True, False), (False, True, False), (True, False, True)]:
    with tempfile.TemporaryDirectory() as tmp:
        app = Path(tmp)
        release = app / "releases" / NEW
        release.mkdir(parents=True)
        shutil.copyfile(SCRIPT, release / "update.sh")
        (release / "portfolio-image.tar.gz").touch()
        (release / "docker-compose.yml").touch()
        (app / ".env").touch()
        if previous:
            old = app / "releases" / OLD
            old.mkdir()
            (old / "docker-compose.yml").touch()
            (app / "current").symlink_to(Path("releases") / OLD)
        binary = app / "bin"
        binary.mkdir()
        docker = binary / "docker"
        docker.write_text('''#!/usr/bin/env python3
import json, os, sys
with open(os.environ["DOCKER_TEST_LOG"], "a") as log:
    log.write(json.dumps([sys.argv[1:], os.environ.get("PORTFOLIO_IMAGE")]) + "\\n")
sys.exit(1 if "up" in sys.argv and os.environ.get("PORTFOLIO_IMAGE") == os.environ.get("FAIL_IMAGE") else 0)
''')
        docker.chmod(0o755)
        log = app / "docker.log"
        env = dict(os.environ, PATH=str(binary) + os.pathsep + os.environ["PATH"],
                   DOCKER_TEST_LOG=str(log), FAIL_IMAGE="puneet-terminal:" + NEW if fail else "")
        with (app / "deploy.lock").open("w") as lock:
            if locked:
                fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
            result = subprocess.run(["bash", str(release / "update.sh")], env=env,
                                    capture_output=True, text=True, timeout=10)
        assert (result.returncode == 0) == (not fail and not locked), result.stderr
        calls = [json.loads(line) for line in log.read_text().splitlines()] if log.exists() else []
        if locked:
            assert not calls, calls
            assert (release / "portfolio-image.tar.gz").exists()
        else:
            assert not (release / "portfolio-image.tar.gz").exists()
            compose_calls = [(args, image) for args, image in calls if args[0] == "compose"]
            for args, _ in compose_calls:
                assert args[1:3] == ["--project-name", "puneet-terminal"], args
                assert args[-1] == "portfolio" or args[-2:] == ["config", "--quiet"], args
                assert not {"down", "prune", "--remove-orphans"}.intersection(args), args
            if fail and previous:
                args, image = compose_calls[-1]
                assert "up" in args and image == "puneet-terminal:" + OLD, compose_calls
            elif fail:
                assert compose_calls[-1][0][-2:] == ["stop", "portfolio"], compose_calls
        if not fail and not locked:
            assert (app / "current").resolve() == release
        elif previous:
            assert (app / "current").resolve() == app / "releases" / OLD
        else:
            assert not (app / "current").exists()
print("Release activation, rollback, first-deploy failure, and locking checks passed.")
