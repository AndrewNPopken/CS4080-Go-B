package main

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/key"
	"image"
	"time"
	//"image/color"
	"image/draw"
	"./camera3d"
	"./space3d"
	"fmt"
	//"math"
)

func fullRender(){

}

func main() {
	ScreenWidth, ScreenHeight := 480, 360
	driver.Main(func(s screen.Screen) {
		nwo := screen.NewWindowOptions{Title:"Window", Height:ScreenHeight, Width:ScreenWidth}
		window, err := s.NewWindow(&nwo)
		if err != nil {
			return
		}
		defer window.Release()
		
		buffer, err := s.NewBuffer(image.Point{nwo.Width, nwo.Height})
		if err != nil {
			return
		}
		defer buffer.Release()
		
		/*
		rgba := buffer.RGBA()
		xMax, yMax := rgba.Rect.Max.X, rgba.Rect.Max.Y
		for x := 0; x < xMax; x++ {
			for y := 0; y < yMax; y++ {
				//rgba.Pix[(y-rgba.Rect.Min.Y)*rgba.Stride + (x-rgba.Rect.Min.X)*4]
//				rgba.Pix[(y)*rgba.Stride + (x)*4]     = 0xff //R
//				rgba.Pix[(y)*rgba.Stride + (x)*4 + 1] = 0x00 //G
//				rgba.Pix[(y)*rgba.Stride + (x)*4 + 2] = 0x00 //B
//				rgba.Pix[(y)*rgba.Stride + (x)*4 + 3] = 0xff //A
				
				var c color.RGBA
				if xd, yd, rMin, rMax := x-240, y-180, 99*99, 10000; xd*xd+yd*yd >= rMin && xd*xd+yd*yd <= rMax {
					c = color.RGBA{0,0,0,0}
				} else {
					c = color.RGBA{uint8( x * 255 / 480 ), uint8(255 - x * 255 / 480 ), uint8( y * 255 / 480 ), 255}
				}
				
				rgba.SetRGBA(x, y, c)
			}
		}
		*/
		
		//create the scene (objects and lights)
		var objects []camera3d.Object
		var lights []camera3d.Light
		var camera camera3d.Camera
		camera.ToWorld = space3d.NewIdentityMatrix()
		//set up options
		//var Width, Height, Depth int
		//var FieldOfView float64
		options := camera3d.Options{Width: ScreenWidth, Height: ScreenHeight, Depth: 0, FieldOfView: 90.0}

		//render
		render := camera3d.Render(&camera, objects, lights, &options)
		
		//draw rendered image onto the screen's buffer's image representation (DO NOT REASSIGN BUFFER)
		rgba := buffer.RGBA()
		draw.Draw(rgba, rgba.Bounds(), render, image.Point{0, 0}, 0)
		
		window.Upload(image.Point{0,0}, buffer, buffer.Bounds())
		window.Publish()
		
		/*//Loopable:
		{
			render = camera3d.Render(&camera, objects, lights, &options)
			draw.Draw(rgba, rgba.Bounds(), render, image.Point{0, 0}, 0)
			window.Upload(image.Point{0,0}, buffer, buffer.Bounds())
			window.Publish()
		}
		*/
		onesecond, _ := time.ParseDuration("1s")
		minFrameRate, _ := time.ParseDuration("0.1s")
		timenext := time.Now().Add(onesecond)
		timelast := time.Now()
		timedelta := time.Since(timelast).Seconds()
		framecount:= -1
		lookup, lookdown, lookleft, lookright := false, false, false, false
		eventChannel := make(chan interface{})
		go func () {
			for {
				eventChannel <- window.NextEvent()
			}
		}()
		for {
			framecount++
				//fmt.Println("FPS: ",framecount)
			if time.Now().After(timenext) {
				fmt.Println("FPS: ",framecount)
				framecount = 0
				timenext = timenext.Add(onesecond)
			}
			
			moreEvents := true
			for moreEvents && !time.Now().After(timelast.Add(minFrameRate)) {
				
				select {
				case ev := <- eventChannel:
					//1 or more events
				
					//fmt.Printf("%T: %v\n", ev, ev)
					switch e := ev.(type) {
					case lifecycle.Event:
						if e.To == lifecycle.StageDead {
							return
						}
					case key.Event:
						//fmt.Println(e.Code)
						if e.Direction == key.DirPress {
							switch e.Code {
							case key.CodeUpArrow:
								lookup = true
							case key.CodeDownArrow:
								lookdown = true
							case key.CodeLeftArrow:
								lookleft = true
							case key.CodeRightArrow:
								lookright = true
							}
						} else if e.Direction == key.DirRelease {
							switch e.Code {
							case key.CodeUpArrow:
								lookup = false
							case key.CodeDownArrow:
								lookdown = false
							case key.CodeLeftArrow:
								lookleft = false
							case key.CodeRightArrow:
								lookright = false
							}
						}
					}
					
				
				default:
					//no event
					moreEvents = false
				}
			}
			if (lookleft && !lookright){
				camera.TurnLeft(1.5 * timedelta)
			} else if (!lookleft && lookright){
				camera.TurnLeft(-1.5 * timedelta)
			}
			if (lookup && !lookdown){
				camera.TurnUp(1.5 * timedelta)
			} else if (!lookup && lookdown){
				camera.TurnUp(-1.5 * timedelta)
			}
			
			timedelta = time.Since(timelast).Seconds()
			timelast = time.Now()
			
			render = camera3d.Render(&camera, objects, lights, &options)
			draw.Draw(rgba, rgba.Bounds(), render, image.Point{0, 0}, 0)
			window.Upload(image.Point{0,0}, buffer, buffer.Bounds())
			window.Publish()
		}
	})
}