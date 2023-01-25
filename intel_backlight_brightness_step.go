package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var brightness_file = "/sys/class/backlight/intel_backlight/brightness"
var max_brightness_file = "/sys/class/backlight/intel_backlight/max_brightness"
var change_percent float32 = 10
var debug = false

func read_int_from_file(filepath string) int {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	s := strings.TrimRight(string(b), "\n")
	number, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return number
}

func max_brightness() int {
	return read_int_from_file(max_brightness_file)
}

func current_brightness() int {
	return read_int_from_file(brightness_file)
}

func change_brightness(val int) {
	data := []byte(strconv.Itoa(val) + "\n")
	err := ioutil.WriteFile(brightness_file, data, 0764)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		brightness := current_brightness()
		maximum_brightness := max_brightness()
		brightness_step := int(float32(maximum_brightness) * 0.01 * change_percent)
		if brightness_step < 2 {
			log.Fatal("calculated brightness step is ", brightness_step)
			return
		}
		if debug {
			log.Println("maximum brightness:", maximum_brightness)
			log.Println("brightness:", brightness)
			log.Println("step:", brightness_step)
		}
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					actual_brightness := current_brightness()
					if debug {
						log.Println("brightness:", actual_brightness)
					}
					diff := actual_brightness - brightness
					if diff == 1 || diff == -1 {
						if debug {
							log.Println("brightness changed by 1")
						}
						if diff == 1 {
							brightness = actual_brightness + brightness_step
							if brightness > maximum_brightness {
								brightness = maximum_brightness
							}
						} else {
							brightness = actual_brightness - brightness_step
							if brightness < 2 {
								brightness = 1
							}
						}
						change_brightness(brightness)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path to watcher
	err = watcher.Add(brightness_file)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
