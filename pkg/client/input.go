package client

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/jeffbaumes/gogame/pkg/client/scene"
	"github.com/jeffbaumes/gogame/pkg/common"
)

func cursorGrabbed(w *glfw.Window) bool {
	return w.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled
}

func cursorPosCallback(player *common.Player) func(w *glfw.Window, xpos, ypos float64) {
	lastCursor := mgl64.Vec2{math.NaN(), math.NaN()}
	return func(w *glfw.Window, xpos, ypos float64) {
		if !cursorGrabbed(w) {
			lastCursor = mgl64.Vec2{math.NaN(), math.NaN()}
			return
		}
		curCursor := mgl64.Vec2{xpos, ypos}
		if math.IsNaN(lastCursor[0]) {
			lastCursor = curCursor
		}
		delta := curCursor.Sub(lastCursor)
		player.Swivel(delta[0], delta[1])
		lastCursor = curCursor
	}
}

func glToPixel(w *glfw.Window, xpos, ypos float64) (xpix, ypix float64) {
	winw, winh := w.GetSize()
	return float64(winw) * (xpos + 1) / 2, float64(winh) * (-ypos + 1) / 2
}
func Sendtext(player *common.Player, text string) {
	if text == "enter" {
		player.Intext = false
		player.Text = ""
		// need to send text to server
	} else if text == "delete" {
		if len(player.Text) != 0 {
			player.Text = player.Text[0:(len(player.Text) - 1)]
		}
	} else {
		player.Text = player.Text + text
	}
}
func keyCallback(player *common.Player) func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if player.InInventory {
			slot := -1
			switch action {
			case glfw.Press:
				switch key {
				case glfw.Key1:
					slot = 0
				case glfw.Key2:
					slot = 1
				case glfw.Key3:
					slot = 2
				case glfw.Key4:
					slot = 3
				case glfw.Key5:
					slot = 4
				case glfw.Key6:
					slot = 5
				case glfw.Key7:
					slot = 6
				case glfw.Key8:
					slot = 7
				case glfw.Key9:
					slot = 8
				case glfw.Key0:
					slot = 9
				case glfw.KeyMinus:
					slot = 10
				case glfw.KeyEqual:
					slot = 11
				}
			}
			if slot >= 0 {
				xpos, ypos := w.GetCursorPos()
				winw, winh := w.GetSize()
				aspect := float32(winw) / float32(winh)
				for m := range common.Materials {
					sz := float32(0.05)
					px := 1.25 * 2 * sz * (float32(m) - float32(len(common.Materials))/2)
					py := 1 - 0.25*aspect
					scale := sz
					xMin, yMin := glToPixel(w, float64(px-scale), float64(py+scale*aspect))
					xMax, yMax := glToPixel(w, float64(px+scale), float64(py-scale*aspect))
					if float64(xpos) >= xMin && float64(xpos) <= xMax && float64(ypos) >= yMin && float64(ypos) <= yMax {
						player.Hotbar[slot] = m
					}
				}
			}
		}
		if player.Intext == false {
			switch action {
			case glfw.Press:
				switch key {
				case glfw.KeySpace:
					if player.GameMode == common.Normal {
						player.HoldingJump = true
					} else {
						player.UpVel = player.WalkVel
					}
				case glfw.KeyLeftShift:
					if player.GameMode == common.Flying {
						player.DownVel = player.WalkVel
					}
				case glfw.Key1:
					player.ActiveHotBarSlot = 0
				case glfw.Key2:
					player.ActiveHotBarSlot = 1
				case glfw.Key3:
					player.ActiveHotBarSlot = 2
				case glfw.Key4:
					player.ActiveHotBarSlot = 3
				case glfw.Key5:
					player.ActiveHotBarSlot = 4
				case glfw.Key6:
					player.ActiveHotBarSlot = 5
				case glfw.Key7:
					player.ActiveHotBarSlot = 6
				case glfw.Key8:
					player.ActiveHotBarSlot = 7
				case glfw.Key9:
					player.ActiveHotBarSlot = 8
				case glfw.Key0:
					player.ActiveHotBarSlot = 9
				case glfw.KeyMinus:
					player.ActiveHotBarSlot = 10
				case glfw.KeyEqual:
					player.ActiveHotBarSlot = 11
				case glfw.KeyE:
					if player.InInventory == false {
						player.HotbarOn = true
						player.InInventory = true
						w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
					} else {
						player.InInventory = false
						w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
					}
				case glfw.KeyT:
					if player.Intext == true {
						player.Intext = false
					} else {
						player.Intext = true
					}
				case glfw.KeyH:
					if player.HotbarOn == false {
						player.HotbarOn = true
					} else {
						player.HotbarOn = false
					}
				case glfw.KeyW:
					player.ForwardVel = player.WalkVel
				case glfw.KeyS:
					player.BackVel = player.WalkVel
				case glfw.KeyD:
					player.RightVel = player.WalkVel
				case glfw.KeyA:
					player.LeftVel = player.WalkVel
				case glfw.KeyM:
					player.GameMode++
					if player.GameMode >= common.NumGameModes {
						player.GameMode = 0
					}
					if player.GameMode == common.Flying {
						player.FallVel = 0
					}
				case glfw.KeyEscape:
					w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
				}
			case glfw.Release:
				switch key {
				case glfw.KeySpace:
					player.HoldingJump = false
					player.UpVel = 0
				case glfw.KeyLeftShift:
					player.DownVel = 0
				case glfw.KeyW:
					player.ForwardVel = 0
				case glfw.KeyS:
					player.BackVel = 0
				case glfw.KeyD:
					player.RightVel = 0
				case glfw.KeyA:
					player.LeftVel = 0
				}
			}
		} else {
			switch action {
			case glfw.Press:
				switch key {
				case glfw.KeyEscape:
					player.Intext = false
				case glfw.Key1:
					Sendtext(player, "1")
				case glfw.Key2:
					Sendtext(player, "2")
				case glfw.Key3:
					Sendtext(player, "3")
				case glfw.Key4:
					Sendtext(player, "4")
				case glfw.Key5:
					Sendtext(player, "5")
				case glfw.Key6:
					Sendtext(player, "6")
				case glfw.Key7:
					Sendtext(player, "7")
				case glfw.Key8:
					Sendtext(player, "8")
				case glfw.Key9:
					Sendtext(player, "9")
				case glfw.Key0:
					Sendtext(player, "0")
				case glfw.KeyQ:
					Sendtext(player, "q")
				case glfw.KeyW:
					Sendtext(player, "w")
				case glfw.KeyE:
					Sendtext(player, "e")
				case glfw.KeyR:
					Sendtext(player, "r")
				case glfw.KeyT:
					Sendtext(player, "t")
				case glfw.KeyY:
					Sendtext(player, "y")
				case glfw.KeyU:
					Sendtext(player, "u")
				case glfw.KeyI:
					Sendtext(player, "i")
				case glfw.KeyO:
					Sendtext(player, "o")
				case glfw.KeyP:
					Sendtext(player, "p")
				case glfw.KeyA:
					Sendtext(player, "a")
				case glfw.KeyS:
					Sendtext(player, "s")
				case glfw.KeyD:
					Sendtext(player, "d")
				case glfw.KeyF:
					Sendtext(player, "f")
				case glfw.KeyG:
					Sendtext(player, "g")
				case glfw.KeyH:
					Sendtext(player, "h")
				case glfw.KeyJ:
					Sendtext(player, "j")
				case glfw.KeyK:
					Sendtext(player, "k")
				case glfw.KeyL:
					Sendtext(player, "l")
				case glfw.KeyZ:
					Sendtext(player, "z")
				case glfw.KeyX:
					Sendtext(player, "x")
				case glfw.KeyC:
					Sendtext(player, "c")
				case glfw.KeyV:
					Sendtext(player, "v")
				case glfw.KeyB:
					Sendtext(player, "b")
				case glfw.KeyN:
					Sendtext(player, "n")
				case glfw.KeyM:
					Sendtext(player, "m")
				case glfw.KeyEnter:
					Sendtext(player, "enter")
				case glfw.KeyDelete:
					Sendtext(player, "delete")
				case glfw.KeyBackspace:
					Sendtext(player, "delete")
				case glfw.KeySpace:
					Sendtext(player, " ")

				}
			}
		}
	}
}

func windowSizeCallback(w *glfw.Window, wd, ht int) {
	fwidth, fheight := scene.FramebufferSize(w)
	gl.Viewport(0, 0, int32(fwidth), int32(fheight))
}

func mouseButtonCallback(player *common.Player, planetRen *scene.Planet) func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	return func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		planet := planetRen.Planet
		if cursorGrabbed(w) {
			if action == glfw.Press && button == glfw.MouseButtonLeft {
				increment := player.LookDir().Mul(0.05)
				pos := player.Loc
				for i := 0; i < 100; i++ {
					pos = pos.Add(increment)
					cell := planet.CartesianToCell(pos)
					if cell != nil && cell.Material != common.Air {
						cellIndex := planet.CartesianToCellIndex(pos)
						planetRen.SetCellMaterial(cellIndex, common.Air)
						break
					}
				}
			} else if action == glfw.Press && button == glfw.MouseButtonRight {
				increment := player.LookDir().Mul(0.05)
				pos := player.Loc
				prevCellIndex := common.CellIndex{Lon: -1, Lat: -1, Alt: -1}
				cellIndex := common.CellIndex{}
				for i := 0; i < 100; i++ {
					pos = pos.Add(increment)
					nextCellIndex := planet.CartesianToCellIndex(pos)
					if nextCellIndex != cellIndex {
						prevCellIndex = cellIndex
						cellIndex = nextCellIndex
						cell := planet.CellIndexToCell(cellIndex)
						if cell != nil && cell.Material != common.Air {
							if prevCellIndex.Lon != -1 {
								planetRen.SetCellMaterial(prevCellIndex, player.Hotbar[player.ActiveHotBarSlot])
							}
							break
						}
					}
				}
			}
		} else {
			if action == glfw.Press {
				w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
			}
		}
	}
}
