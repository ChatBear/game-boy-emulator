import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import EmulatorArchitecture from './EmulationScreen'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <EmulatorArchitecture />
  </StrictMode>,
)
