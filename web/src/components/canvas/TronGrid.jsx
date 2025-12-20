import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'

export default function TronGrid({ size = 100, divisions = 50 }) {
  const linesRef = useRef()

  const linePositions = useMemo(() => {
    const positions = []
    const half = size / 2
    const step = size / divisions

    // Horizontal lines
    for (let i = 0; i <= divisions; i++) {
      const z = -half + i * step
      positions.push(-half, -5, z, half, -5, z)
    }

    // Vertical lines  
    for (let i = 0; i <= divisions; i++) {
      const x = -half + i * step
      positions.push(x, -5, -half, x, -5, half)
    }

    return new Float32Array(positions)
  }, [size, divisions])

  useFrame((state) => {
    if (linesRef.current) {
      // Pulse effect
      const pulse = Math.sin(state.clock.elapsedTime * 2) * 0.1 + 0.9
      linesRef.current.material.opacity = 0.3 * pulse
    }
  })

  return (
    <group>
      {/* Grid lines */}
      <lineSegments ref={linesRef}>
        <bufferGeometry>
          <bufferAttribute
            attach="attributes-position"
            count={linePositions.length / 3}
            array={linePositions}
            itemSize={3}
          />
        </bufferGeometry>
        <lineBasicMaterial
          color="#333333"
          transparent
          opacity={0.3}
          blending={THREE.AdditiveBlending}
        />
      </lineSegments>

      {/* Glowing center line */}
      <mesh position={[0, -4.9, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[0.5, size]} />
        <meshBasicMaterial
          color="#FFC837"
          transparent
          opacity={0.5}
          blending={THREE.AdditiveBlending}
        />
      </mesh>
    </group>
  )
}
