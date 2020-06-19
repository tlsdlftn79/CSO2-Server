package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func onGameStart(seq *uint8, p packet, client net.Conn) {
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to start game but not in server !")
		return
	}
	//检查用户是不是房主
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to start game in a null room !")
		return
	}
	//房主开始游戏,设置房间状态
	u := rm.roomGetUser(uPtr.userid)
	if u == nil ||
		u.userid <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to start game but is null in room !")
		return
	}
	if rm.hostUserID == uPtr.userid {
		rm.stopCountdown()
		rm.setStatus(StatusIngame)
		//设置用户状态
		uPtr.setUserIngame(true)
		u.setUserIngame(true)
		//对非房主用户发送数据包
		for _, v := range rm.users {
			if v.userid != u.userid {
				rst := BytesCombine(BuildHeader(v.currentSequence, p), buildRoomSetting(*rm))
				sendPacket(rst, v.currentConnection)
				if v.isUserReady() {
					v.setUserIngame(true)
					//连接到主机
					rst = UDPBuild(v.currentSequence, p, 1, u.userid, u.netInfo.ExternalIpAddress, u.netInfo.ExternalServerPort)
					sendPacket(rst, v.currentConnection)
					//加入主机
					p.id = TypeHost
					rst = BytesCombine(BuildHeader(v.currentSequence, p), BuildJoinHost(u.userid))
					sendPacket(rst, v.currentConnection)
					//给主机发送其他人的数据
					rst = UDPBuild(seq, p, 0, v.userid, v.netInfo.ExternalIpAddress, v.netInfo.ExternalClientPort)
					sendPacket(rst, uPtr.currentConnection)
				}
			}
		}
		//给每个人发送房间内所有人的准备状态
		p.id = TypeRoom
		for _, v := range rm.users {
			rst := BuildUserReadyStatus(v)
			for _, k := range rm.users {
				rst = BytesCombine(BuildHeader(k.currentSequence, p), rst)
				sendPacket(rst, k.currentConnection)
			}
		}
		//主机开始游戏
		p.id = TypeHost
		rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildGameStart(u.userid))
		sendPacket(rst, uPtr.currentConnection)
		log.Println("Host", string(uPtr.username), "started game ")
	} else if rm.setting.isIngame != 0 {
		host := rm.roomGetUser(rm.hostUserID)
		if host == nil ||
			host.userid <= 0 {
			log.Println("Error : User", string(uPtr.username), "try to start game but host is null !")
			return
		}
		//发送房间数据
		rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), buildRoomSetting(*rm))
		sendPacket(rst, uPtr.currentConnection)
		//设置用户状态
		uPtr.setUserIngame(true)
		u.setUserIngame(true)
		//连接到主机
		rst = UDPBuild(u.currentSequence, p, 1, host.userid, host.netInfo.ExternalIpAddress, host.netInfo.ExternalServerPort)
		sendPacket(rst, u.currentConnection)
		//加入主机
		p.id = TypeHost
		rst = BytesCombine(BuildHeader(u.currentSequence, p), BuildJoinHost(host.userid))
		sendPacket(rst, u.currentConnection)
		//发送用户数据给主机
		rst = UDPBuild(host.currentSequence, p, 0, uPtr.userid, uPtr.netInfo.ExternalIpAddress, uPtr.netInfo.ExternalServerPort)
		sendPacket(rst, host.currentConnection)
		//给每个人发送房间内所有人的准备状态
		p.id = TypeRoom
		for _, v := range rm.users {
			rst = BuildUserReadyStatus(v)
			for _, k := range rm.users {
				rst = BytesCombine(BuildHeader(k.currentSequence, p), rst)
				sendPacket(rst, k.currentConnection)
			}
		}
		log.Println("User", string(uPtr.username), "started game ")
	}
}

func BuildJoinHost(id uint32) []byte {
	buf := make([]byte, 13)
	offset := 0
	WriteUint8(&buf, HostJoin, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint64(&buf, 0, &offset)
	return buf[:offset]
}

func BuildGameStart(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, GameStart, &offset)
	WriteUint32(&buf, id, &offset)
	return buf[:offset]
}