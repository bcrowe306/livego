package rtmp

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"livego/protocol/rtmp/rtmprelay"
	"livego/router"

	"github.com/gwuhaolin/livego/av"
	"github.com/gwuhaolin/livego/protocol/rtmp/cache"
	"github.com/orcaman/concurrent-map"
)

var (
	EmptyID = ""
)

type RtmpStream struct {
	streams cmap.ConcurrentMap //key
}

func NewRtmpStream() *RtmpStream {
	ret := &RtmpStream{
		streams: cmap.New(),
	}
	go ret.CheckAlive()
	return ret
}

func (rs *RtmpStream) HandleReader(r av.ReadCloser) {
	info := r.Info()
	log.Printf("HandleReader: info[%v]", info)

	var stream *Stream
	i, ok := rs.streams.Get(info.Key)
	if stream, ok = i.(*Stream); ok {
		stream.TransStop()
		id := stream.ID()
		if id != EmptyID && id != info.UID {
			ns := NewStream()
			stream.Copy(ns)
			stream = ns
			rs.streams.Set(info.Key, ns)
		}
	} else {
		stream = NewStream()
		rs.streams.Set(info.Key, stream)
		stream.info = info
	}

	stream.AddReader(r)
}

func (rs *RtmpStream) HandleWriter(w av.WriteCloser) {
	info := w.Info()
	log.Printf("HandleWriter: info[%v]", info)

	// TODO: Add logic to deny live playback if
	log.Println(info.Key, info.URL)
	var s *Stream
	ok := rs.streams.Has(info.Key)
	if !ok {
		s = NewStream()
		rs.streams.Set(info.Key, s)
		s.info = info
	} else {
		item, ok := rs.streams.Get(info.Key)
		if ok {
			s = item.(*Stream)
			s.AddWriter(w)
		}
	}
}

func (rs *RtmpStream) GetStreams() cmap.ConcurrentMap {
	return rs.streams
}

func (rs *RtmpStream) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		for item := range rs.streams.IterBuffered() {
			v := item.Val.(*Stream)
			if v.CheckAlive() == 0 {
				rs.streams.Remove(item.Key)
			}
		}
	}
}

type Stream struct {
	isStart bool
	cache   *cache.Cache
	r       av.ReadCloser
	ws      cmap.ConcurrentMap
	info    av.Info
}

type PackWriterCloser struct {
	init bool
	w    av.WriteCloser
}

func (p *PackWriterCloser) GetWriter() av.WriteCloser {
	return p.w
}

func NewStream() *Stream {
	return &Stream{
		cache: cache.NewCache(),
		ws:    cmap.New(),
	}
}

func (s *Stream) ID() string {
	if s.r != nil {
		return s.r.Info().UID
	}
	return EmptyID
}

func (s *Stream) GetReader() av.ReadCloser {
	return s.r
}

func (s *Stream) GetWs() cmap.ConcurrentMap {
	return s.ws
}

func (s *Stream) Copy(dst *Stream) {
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		s.ws.Remove(item.Key)
		v.w.CalcBaseTimestamp()
		dst.AddWriter(v.w)
	}
}

func (s *Stream) AddReader(r av.ReadCloser) {
	s.r = r
	go s.TransStart()
}

func (s *Stream) AddWriter(w av.WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	s.ws.Set(info.UID, pw)
}

/*Check if static_push is configured in this application.
If configured, start the push remote connection */
func (s *Stream) StartStaticPush() {
	key := s.info.Key

	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	// This variable is used to add the streamkey to the end of the url. I do not want this behavior to take place here
	// streamname := key[index+1:]
	appname := dscr[0]

	log.Printf("StartStaticPush: current， appname=%s", appname)
	r, exist := router.SelectRoute(key)

	if !exist {
		log.Printf("StartStaticPush: Route not found for url: %v", key)
	}

	// Loop over Route Endpoints, creating push urls
	for _, e := range r.Endpoints {
		if e.Enabled {
			var pushurl string
			if r.CopyKey {
				pushurl = fmt.Sprintf("rtmp://%v/%v/%v", e.Host, e.App, key)
			} else {
				pushurl = fmt.Sprintf("rtmp://%v/%v/%v", e.Host, e.App, e.Stream)
			}
			// Add endpoint to Push Map Object and retrieve
			staticpushObj := rtmprelay.GetAndCreateStaticPushObject(pushurl)
			if staticpushObj != nil {
				if err := staticpushObj.Start(); err != nil {
					log.Printf("StartStaticPush: staticpushObj.Start %s error=%v", pushurl, err)
				} else {
					log.Printf("StartStaticPush: staticpushObj.Start %s ok", pushurl)
				}
			} else {
				log.Printf("StartStaticPush GetStaticPushObject %s error", pushurl)
			}
			log.Printf("StartStaticPush: static pushurl=%s", pushurl)
		}
	}
}

func (s *Stream) StopStaticPush() {
	key := s.info.Key

	log.Printf("StopStaticPush......%s", key)
	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	streamname := key[index+1:]
	appname := dscr[0]

	log.Printf("StopStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist, err := rtmprelay.GetStaticPushList(key)
	if err != nil || len(pushurllist) < 1 {
		log.Printf("StopStaticPush: GetStaticPushList error=%v", err)
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		log.Printf("StopStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := rtmprelay.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.Stop()
			rtmprelay.ReleaseStaticPushObject(pushurl)
			log.Printf("StopStaticPush: staticpushObj.Stop %s ", pushurl)
		} else {
			log.Printf("StopStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (s *Stream) IsSendStaticPush() bool {
	key := s.info.Key
	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return false
	}

	//log.Printf("SendStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist, err := rtmprelay.GetStaticPushList(key)
	if err != nil || len(pushurllist) < 1 {
		//log.Printf("SendStaticPush: GetStaticPushList error=%v", err)
		return false
	}

	// Check stream URL for correct format
	index := strings.Index(key, "/")
	if index < 0 {
		return false
	}

	// Enumertate pushlist
	for _, pushurl := range pushurllist {
		//log.Printf("SendStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := rtmprelay.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			return true
			//staticpushObj.WriteAvPacket(&packet)
			//log.Printf("SendStaticPush: WriteAvPacket %s ", pushurl)
		} else {
			log.Printf("SendStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
	return false
}

func (s *Stream) SendStaticPush(packet av.Packet) {
	key := s.info.Key

	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	// This variable was used to place the streamkey on the end of the URL. I don't want this behavior here.
	// streamname := key[index+1:]
	// appname := dscr[0]

	//log.Printf("SendStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist, err := rtmprelay.GetStaticPushList(key)
	if err != nil || len(pushurllist) < 1 {
		//log.Printf("SendStaticPush: GetStaticPushList error=%v", err)
		return
	}

	for _, pushurl := range pushurllist {

		// This code was taking the streamkey and placing it on the end of the url. I dont want this behavior to take place here
		// pushurl := pushurl + "/" + streamname
		//log.Printf("SendStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := rtmprelay.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.WriteAvPacket(&packet)
			//log.Printf("SendStaticPush: WriteAvPacket %s ", pushurl)
		} else {
			log.Printf("SendStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (s *Stream) TransStart() {
	s.isStart = true
	var p av.Packet

	log.Printf("TransStart:%v", s.info)

	s.StartStaticPush()

	for {
		if !s.isStart {
			s.closeInter()
			return
		}
		err := s.r.Read(&p)
		if err != nil {
			s.closeInter()
			s.isStart = false
			return
		}

		if s.IsSendStaticPush() {
			s.SendStaticPush(p)
		}

		s.cache.Write(p)

		for item := range s.ws.IterBuffered() {
			v := item.Val.(*PackWriterCloser)
			if !v.init {
				//log.Printf("cache.send: %v", v.w.Info())
				if err = s.cache.Send(v.w); err != nil {
					log.Printf("[%s] send cache packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
					continue
				}
				v.init = true
			} else {
				new_packet := p
				//writeType := reflect.TypeOf(v.w)
				//log.Printf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&new_packet); err != nil {
					log.Printf("[%s] write packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
				}
			}
		}
	}
}

func (s *Stream) TransStop() {
	log.Printf("TransStop: %s", s.info.Key)

	if s.isStart && s.r != nil {
		s.r.Close(errors.New("stop old"))
	}

	s.isStart = false
}

func (s *Stream) CheckAlive() (n int) {
	if s.r != nil && s.isStart {
		if s.r.Alive() {
			n++
		} else {
			s.r.Close(errors.New("read timeout"))
		}
	}
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			if !v.w.Alive() && s.isStart {
				s.ws.Remove(item.Key)
				v.w.Close(errors.New("write timeout"))
				continue
			}
			n++
		}

	}
	return
}

func (s *Stream) closeInter() {
	if s.r != nil {
		s.StopStaticPush()
		log.Printf("[%v] publisher closed", s.r.Info())
	}

	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			if v.w.Info().IsInterval() {
				v.w.Close(errors.New("closed"))
				s.ws.Remove(item.Key)
				log.Printf("[%v] player closed and remove\n", v.w.Info())
			}
		}

	}
}
