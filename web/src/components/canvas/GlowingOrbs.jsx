import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { Float, Sphere } from '@react-three/drei'
import * as THREE from 'three'

export default function GlowingOrbs({ count = 8 }) {
  const groupRef = useRef()

  useFrame((state) => {
    if (groupRef.current) {
      groupRef.current.rotation.y = state.clock.elapsedTime * 0.1
    }
  })

  const orbs = Array.from({ length: count }, (_, i) => {
    const angle = (i / count) * Math.PI * 2
    const radius = 8 + Math.random() * 4
    const y = (Math.random() - 0.5) * 10
    
    return {
      position: [
        Math.cos(angle) * radius,
        y,
        Math.sin(angle) * radius
      ],
      scale: 0.1 + Math.random() * 0.3,
      color: i % 2 === 0 ? '#00d4ff' : '#0066ff',
      speed: 0.5 + Math.random() * 1
    }
  })

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
