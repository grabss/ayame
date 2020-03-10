package main

import (
	"regexp"
	"strings"
)

var (
	// register/unregister は待たせる
	registerChannel   = make(chan *register)
	unregisterChannel = make(chan *unregister)
	// ブロックされたくないので 100 に設定
	forwardChannel = make(chan forward, 100)
	knockChannel = make(chan knock, 100)
)

// roomId がキーになる
type room struct {
	connections map[string]*connection
}

func server() {
	// room を管理するマップはここに用意する
	var m = make(map[string]room)
	// ここはシングルなのでロックは不要、多分
	for {
		select {
		case register := <-registerChannel:
			c := register.connection
			rch := register.resultChannel
			r, ok := m[c.roomID]
			if ok {
				// room があった
				if len(r.connections) == 1 {
					r.connections[c.ID] = c
					m[c.roomID] = r
					rch <- two
				} else {
					// room あったけど満杯
					rch <- full
				}
			} else {
				// room がなかった
				var connections = make(map[string]*connection)
				connections[c.ID] = c
				// room を追加
				m[c.roomID] = room{
					connections: connections,
				}
				c.debugLog().Msg("CREATED-ROOM")
				rch <- one
			}
		case unregister := <-unregisterChannel:
			c := unregister.connection
			// room を探す
			r, ok := m[c.roomID]
			// room がない場合は何もしない
			if ok {
				_, ok := r.connections[c.ID]
				if ok {
					for _, connection := range r.connections {
						// 両方の forwardChannel を閉じる
						close(connection.forwardChannel)
						connection.debugLog().Msg("CLOSED-FORWARD-CHANNEL")
						connection.debugLog().Msg("REMOVED-CLIENT")
					}
					// room を削除
					delete(m, c.roomID)
					c.debugLog().Msg("DELETED-ROOM")
				}
			}
		case forward := <-forwardChannel:
			r, ok := m[forward.connection.roomID]
			// room がない場合は何もしない
			if ok {
				// room があった
				for connectionID, client := range r.connections {
					// 自分ではない方に投げつける
					if connectionID != forward.connection.ID {
						client.forwardChannel <- forward
					}
				}
			}
		case knock := <-knockChannel:
			rch := knock.resultChannel
			var result = false
			if rm, _ := regexp.MatchString("^\\d{1,}-\\d{1,}$", knock.knockID); rm {
				// roomID検索の場合(d-d)
				// -> 純粋なIndex検索で判定
				_, result = m[knock.knockID]
			} else if rm, _ := regexp.MatchString("^\\d{1,}$", knock.knockID); rm {
				// clientID検索の場合(d)
				// -> 登録済みのroomIDを回しd-dのみ検索対象にした上でハイフンでsplitし両方のclientIDで判定
				loopRooms:
					for k, _ := range m {
						if km, _ := regexp.MatchString("^\\d{1,}-\\d{1,}$", k); km {
							for _, v := range strings.Split(k, "-") {
								if v == knock.knockID {
									result = true
									break loopRooms
								}
							}
						}
					}
			}
			rch <- result
		}
	}
}
