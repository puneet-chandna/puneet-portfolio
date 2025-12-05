import { Suspense, useRef, useEffect } from 'react'
import { Canvas, useFrame, useThree } from '@react-three/fiber'
import { PerspectiveCamera } from '@react-three/drei'
import { EffectComposer, Bloom, ChromaticAberration } from '@react-three/postprocessing'
import { BlendFunction } from 'postprocessing'
import ParticleField from './ParticleField'
import TronGrid from './TronGrid'
import GlowingOrbs from './GlowingOrbs'

// Mouse parallax camera controller
function CameraController() {
  const { camera } = useThree()
  const mouseRef = useRef({ x: 0, y: 0 })
  const targetRef = useRef({ x: 0, y: 0 })

  useEffect(() => {
    const handleMouseMove = (e) => {
      mouseRef.current.x = (e.clientX / window.innerWidth - 0.5) * 2
      mouseRef.current.y = (e.clientY / window.innerHeight - 0.5) * 2
    }
    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  useFrame(() => {
    // Smooth lerp towards mouse position
    targetRef.current.x += (mouseRef.current.x - targetRef.current.x) * 0.05
    targetRef.current.y += (mouseRef.current.y - targetRef.current.y) * 0.05

    // Apply parallax to camera position using set()
    camera.position.set(
      targetRef.current.x * 3,
      5 - targetRef.current.y * 2,
      20
    )
    camera.lookAt(0, 0, 0)
  })

  return null
}

export default function Scene() {
  return (
    <div className="canvas-container">
      <Canvas>
        <PerspectiveCamera makeDefault position={[0, 5, 20]} fov={60} />
        <color attach="background" args={['#0a0a0f']} />
        <fog attach="fog" args={['#0a0a0f', 20, 60]} />
        
        <ambientLight intensity={0.2} />
        <pointLight position={[10, 10, 10]} intensity={0.5} color="#00d4ff" />
        
        <Suspense fallback={null}>
          <ParticleField count={1500} />
          <TronGrid size={80} divisions={40} />
          <GlowingOrbs count={6} />
        </Suspense>

        <CameraController />

        <EffectComposer>
          <Bloom
            intensity={1.5}
            luminanceThreshold={0.1}
            luminanceSmoothing={0.9}
            blendFunction={BlendFunction.ADD}
          />
          <ChromaticAberration
            offset={[0.0005, 0.0005]}
            blendFunction={BlendFunction.NORMAL}
          />
        </EffectComposer>
      </Canvas>
    </div>
  )
}
