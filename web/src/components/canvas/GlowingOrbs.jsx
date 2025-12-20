import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import { Float, Sphere } from '@react-three/drei'
import * as THREE from 'three'

// Seeded random for deterministic results
function seededRandom(seed) {
  const x = Math.sin(seed * 9999) * 10000
  return x - Math.floor(x)
}

export default function GlowingOrbs({ count = 8 }) {
  const groupRef = useRef()

  useFrame((state) => {
    if (groupRef.current) {
      groupRef.current.rotation.y = state.clock.elapsedTime * 0.1
    }
  })

  const orbs = useMemo(() => {
    return Array.from({ length: count }, (_, i) => {
      const angle = (i / count) * Math.PI * 2
      const radius = 8 + seededRandom(i * 10) * 4
      const y = (seededRandom(i * 20) - 0.5) * 10
      
      return {
        position: [
          Math.cos(angle) * radius,
          y,
          Math.sin(angle) * radius
        ],
        scale: 0.1 + seededRandom(i * 30) * 0.3,
        color: i % 3 === 0 ? '#FFC837' : '#E0E0E0',
        speed: 0.5 + seededRandom(i * 40) * 1
      }
    })
  }, [count])

  return (
    <group ref={groupRef}>
      {orbs.map((orb, i) => (
        <Float
          key={i}
          speed={orb.speed}
          rotationIntensity={0.5}
          floatIntensity={2}
        >
          <Sphere
            args={[orb.scale, 16, 16]}
            position={orb.position}
          >
            <meshBasicMaterial
              color={orb.color}
              transparent
              opacity={0.8}
              blending={THREE.AdditiveBlending}
            />
          </Sphere>
          {/* Glow effect */}
          <Sphere
            args={[orb.scale * 2, 16, 16]}
            position={orb.position}
          >
            <meshBasicMaterial
              color={orb.color}
              transparent
              opacity={0.15}
              blending={THREE.AdditiveBlending}
            />
          </Sphere>
        </Float>
      ))}
    </group>
  )
}
