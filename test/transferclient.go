/*
 * 文件传送客户端 
 * 
 */ 

package main

import(
	"net"
	tcpsession "kendynet-go/tcpsession"
	packet "kendynet-go/packet"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	request_file = 1
	file_size = 2
	transfering = 3
)

type transfer_session struct{
	filecontent []byte
	widx        int
	filename    string
	filesize    int
}

func (this *transfer_session)recv_file(rpk *packet.Rpacket)(bool){
	content,_ := rpk.Binary()
	copy(this.filecontent[this.widx:],content[:])
	this.widx += len(content)
	if this.widx >= this.filesize {
		ioutil.WriteFile(this.filename, this.filecontent, 0x777)
		return true
	}
	return false
}

func process_client(session *tcpsession.Tcpsession,rpk *packet.Rpacket){
	cmd,_ := rpk.Uint16()
	if cmd == file_size {
		if session.Ud() == nil {
			fmt.Printf("error\n")
			session.Close()
			return
		}
		tsession := session.Ud().(*transfer_session)
		filesize,_ := rpk.Uint32()
		fmt.Printf("file size:%d\n",filesize)
		tsession.widx = 0
		tsession.filesize = int(filesize)
		tsession.filecontent = make([]byte,filesize)
		
	}else if cmd == transfering {
		if session.Ud() == nil {
			fmt.Printf("close here\n")
			session.Close()
			return
		}
		tsession := session.Ud().(*transfer_session)
		if tsession.recv_file(rpk) {
			//传输完毕
			fmt.Printf("transfer finish\n")
			session.Close()
			return
		}
	}else{
		fmt.Printf("cmd error,%d\n",cmd)
		//session.Close()
	}
}

func session_close(session *tcpsession.Tcpsession){
	fmt.Printf("client disconnect\n")
}

func main(){
	
	if len(os.Args) < 3 {
		fmt.Printf("usage ./transferclient <filename> <savefilename\n")
		return
	}
	service := "127.0.0.1:8010"
	tcpAddr,err := net.ResolveTCPAddr("tcp4", service)
	if err != nil{
		fmt.Printf("ResolveTCPAddr")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("DialTcp error,%s\n",err)
	}else{
		session := tcpsession.NewTcpSession(conn,false)
		fmt.Printf("connect sucessful\n")
		//发出文件请求
		wpk := packet.NewWpacket(packet.NewByteBuffer(64),false)
		wpk.PutUint16(request_file)
		wpk.PutString(os.Args[1])
		session.Send(wpk,nil)
		tsession := &transfer_session{filename:os.Args[2]}
		session.SetUd(tsession)	
		tcpsession.ProcessSession(session,process_client,session_close)
	}
}

