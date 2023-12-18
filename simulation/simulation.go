package simulation

import "math"

// Wave represents the properties of a wave
type Wave struct {
	Frequency  float64 // Frequency in Hz
	Amplitude  float64 // Amplitude of the wave
	Speed      float64 // Speed of sound in the medium (m/s)
	Wavelength float64 // Wavelength (m)
}

// NewWave creates a new Wave with given properties
func NewWave(frequency, amplitude, speed float64) *Wave {
	wavelength := speed / frequency
	return &Wave{
		Frequency:  frequency,
		Amplitude:  amplitude,
		Speed:      speed,
		Wavelength: wavelength,
	}
}

// WaveNumber returns the wave number (k) of the wave
func (w *Wave) WaveNumber() float64 {
	return 2 * math.Pi / w.Wavelength
}

// AngularFrequency returns the angular frequency (ω) of the wave
func (w *Wave) AngularFrequency() float64 {
	return 2 * math.Pi * w.Frequency
}

// Displacement calculates the displacement of the wave at a given point and time
func (w *Wave) Displacement(x, y, z, t float64) float64 {
	// Assuming a propagation along the x-axis for simplicity
	k := w.WaveNumber()
	ω := w.AngularFrequency()
	return w.Amplitude * math.Sin(k*x-ω*t)
}
