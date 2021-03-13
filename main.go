package main

import (
	"fmt"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"log"
	//"github.com/nareix/joy4/av/transcode"
	//"github.com/nareix/joy4/cgo/ffmpeg"
	"github.com/nareix/joy4/format"
	"time"
)
func init() {
	format.RegisterAll()
}

func main() {
	outfile, _ := avutil.Create("out.mp4")

	fmt.Println("Start recording...")

	stream([]string{"1.flv"},outfile)

	if err := outfile.WriteTrailer(); err != nil {
		return
	}

	outfile.Close()
	//t1 := time.Now()
	//fmt.Printf("Stop recording, took %v to run.\n", t1.Sub(t0))
}


func stream(paths  []string, outfile av.MuxCloser) {
	files := make(map[string]av.DemuxCloser)
	for _,path := range paths {
		infile, _ := avutil.Open(path)
		files[path] = infile
	}
	streams,_ := files[paths[0]].Streams()
	if err := outfile.WriteHeader(streams); err != nil {
		return
	}
	start := false
	var err error
	count := 0
	var totalTime time.Duration
	t0 := time.Now()
	for {
		var pck av.Packet

		if pck, err = files[paths[count]].ReadPacket(); err != nil {
			log.Println("final",err)
			if count < len(paths)-1 {
				count += 1
				continue
			}else{
				break
			}
		}

		if pck.IsKeyFrame {
			start = true
		}
		if !start {
			continue
		}
		if start {
			pck.Time = totalTime
			outfile.WritePacket(pck)
			totalTime += 40 * time.Millisecond
		}
	}
	t1 := time.Now()
	fmt.Printf("Stop recording, took %v to run.\n", t1.Sub(t0))
	for _,path := range paths {
		files[path].Close()
	}
}