import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'

// Seeded random for deterministic results
function seededRandom(seed) {
  const x = Math.sin(seed * 9999) * 10000
  return x - Math.floor(x)
}

export default function ParticleField({ count = 2000 }) {
  const mesh = useRef()
  
  const particles = useMemo(() => {
    const positions = new Float32Array(count * 3)
    const colors = new Float32Array(count * 3)
    
    for (let i = 0; i < count; i++) {
      positions[i * 3] = (seededRandom(i * 3) - 0.5) * 50
      positions[i * 3 + 1] = (seededRandom(i * 3 + 1) - 0.5) * 50
      positions[i * 3 + 2] = (seededRandom(i * 3 + 2) - 0.5) * 50
      
      // Tron colors: cyan to blue gradient
      const t = seededRandom(i * 7)
      colors[i * 3] = t * 0 + (1 - t) * 0
      colors[i * 3 + 1] = t * 0.83 + (1 - t) * 0.4
      colors[i * 3 + 2] = t * 1 + (1 - t) * 1
    }
    
    return { positions, colors }
  }, [count])

  useFrame((state) => {
    if (mesh.current) {
      mesh.current.rotation.y = state.clock.elapsedTime * 0.02
      mesh.current.rotation.x = Math.sin(state.clock.elapsedTime * 0.1) * 0.1
    }
  })

  return (
    <points ref={mesh}>
      <bufferGeometry>
        <bufferAttribute
          attach="attributes-position"
          count={particles.positions.length / 3}
          array={particles.positions}
          itemSize={3}
        />
        <bufferAttribute
          attach="attributes-color"
          count={particles.colors.length / 3}
          array={particles.colors}
          itemSize={3}
        />
      </bufferGeometry>
      <pointsMaterial
        size={0.08}
        vertexColors
        transparent
        opacity={0.8}
        sizeAttenuation
        blending={THREE.AdditiveBlending}
      />
    </points>
  )
}
