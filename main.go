package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kkdai/youtube/v2"
	_ "github.com/mattn/go-sqlite3"
	qrPng "github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Contoh: Unduh video dengan ID "QAGDGja7kbs"
	// youtube.Video.
	// ExampleClient("MoN9ql6Yymw")
	// ExamplePlaylist()
	// err := youtubedr.Download("aBufJ_TcUAw")
	// if err != nil {
	// 	fmt.Println("Gagal mengunduh video:", err)
	// } else {
	// 	fmt.Println("Video berhasil diunduh!")
	// }

	// ig content download
	JalankanWaBot()
	// url := "https://www.instagram.com/p/C-AehoeRyYW/" // Ganti dengan URL Instagram yang ingin Anda periksa
	// // url := "https://www.instagram.com/p/ABC123/" // Ganti dengan URL Instagram yang ingin Anda periksa

	// contentType, err := getContentType(url)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// if contentType == "image/jpeg" {
	// 	fmt.Println("Konten adalah foto (JPEG)")
	// } else if contentType == "video/mp4" {
	// 	fmt.Println("Konten adalah video (MP4)")
	// } else {
	// 	fmt.Println("Konten tidak dikenali")
	// }
}

func ExampleClient(coba string) {
	videoID := coba
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}
	fmt.Println(video.Description)
	fmt.Println(video.Title)
	fmt.Println(video.CaptionTracks)
	fmt.Println(video.PublishDate)
	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
}

// Example usage for playlists: downloading and checking information.
func ExamplePlaylist() {
	playlistID := "PLQZgI7en5XEgM0L1_ZcKmEzxW1sCOVZwP"
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		panic(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	println(header)
	println(strings.Repeat("=", len(header)) + "\n")

	for k, v := range playlist.Videos {
		fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
	}

	/* ----- Downloading the 1st video ----- */
	entry := playlist.Videos[0]
	video, err := client.VideoFromPlaylistEntry(entry)
	if err != nil {
		panic(err)
	}
	// Now it's fully loaded.

	fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create("video.mp4")

	if err != nil {
		panic(err)
	}

	defer file.Close()
	_, err = io.Copy(file, stream)

	if err != nil {
		panic(err)
	}

	println("Downloaded /video.mp4")
}

func getContentType(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	return contentType, nil
}

// func main() {

// }

var client *whatsmeow.Client

func JalankanWaBot() {

	store.DeviceProps.Os = proto.String("MisterKongBot")
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println("erro1")
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {

		fmt.Println("erro2")
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(handleBot)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			fmt.Println("erro3")
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				// qrterminal.Generate(evt.Code, qrterminal.M, os.Stdout)
				qrPng.WriteFile(evt.Code, qrPng.High, 256, `tesqr.png`)
			} else {
				fmt.Println("Login event:", evt.Event)

			}
		}
		// tampungClient["msKong"] = client
	} else {
		// Already logged in, just connect
		err = client.Connect()
		// tampungClient["msKong"] = client
		// 		dt2 := make(mapB4a)
		// 		// dt2 := readJSON("UserClient.json")
		// WriteJsonUser(dt2,"coba1",client)
		// 		fmt.Println("jalan di bawah connectttttttttt")
		if err != nil {
			fmt.Println("erro5")
			panic(err)
		}
		// var v *events.Message
		// v.Info.Sender.User = "6282266353193"
		// fmt.Println("=======================")
		// fmt.Println(v.Info.Sender)
		// fmt.Println("=======================")
		// 		client.SendMessage(context.Background(),v.Info.Sender,&waProto.Message{
		// 			Conversation: proto.String("Kode sudah tidak dapat di gunakan harap melakukan request lagi di aplikasi kongrider"),
		// 			})
		// 		return
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func handleBot(evt interface{}) {
	switch v := evt.(type) {

	case *events.Message:
		// fmt.Println(v.Info.Chat)
		// // jenis jid :
		// // 1. @s.whatsapp.net
		// // 2. @g.us
		fmt.Print()
		fmt.Println("v.Message : ", v.Message)

		fmt.Println("v : ", v)
		fmt.Println("v.Info.Sender : ", v.Info.Sender)
		fmt.Println("v.Info.Chat : ", v.Info.Chat)
		fmt.Println("v.Info.ID : ", v.Info.ID)

		// v.
		// fmt.Println("*v.Message.ExtendedTextMessage.ContextInfo.StanzaID : ", *v.Message.ExtendedTextMessage.ContextInfo.StanzaID)
		// return
		// v.Message :  extendedTextMessage:{
		// text:"Okeee"
		//  previewType:NONE
		// contextInfo:
		// 	{stanzaID:"5D68D85938023F7A5118DD9A5C78BC48"
		// 	participant:"6282266353193@s.whatsapp.net"
		// 	 quotedMessage:{
		// 		conversation:"Coba"
		// 		}
		// 		} inviteLinkGroupTypeV2:DEFAULT}
		// v :  &{{{120363302887276471@g.us 628889827306@s.whatsapp.net true true } 63344F456306D390B9E980E85CF0BB89 0 text ~ 2024-07-31 22:09:02 +0800 +08  false   {  0001-01-01 00:00:00 +0000 UTC} { } <nil> <nil>} extendedTextMessage:{text:"Okeee" previewType:NONE
		// contextInfo:{stanzaID:"5D68D85938023F7A5118DD9A5C78BC48" participant:"6282266353193@s.whatsapp.net" quotedMessage:{conversation:"Coba"}} inviteLinkGroupTypeV2:DEFAULT} false false false false false false false <nil>  0 <nil> extendedTextMessage:{text:"Okeee" previewType:NONE contextInfo:{stanzaID:"5D68D85938023F7A5118DD9A5C78BC48" participant:"6282266353193@s.whatsapp.net" quotedMessage:{conversation:"Coba"}} inviteLinkGroupTypeV2:DEFAULT}}
		oke, err := client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{Text: proto.String("ini balesan"),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: &waE2E.Message{
						Conversation: proto.String(v.Message.GetConversation())},
					// StanzaID: proto.String("A8F303FE9F440C7F4C63CC47FD2D34A1"),
					StanzaID: &v.Info.ID,
					// Participant: proto.String("6282266353193@s.whatsapp.net"),
					Participant: proto.String(fmt.Sprint(v.Info.Sender)),
				},
				InviteLinkGroupTypeV2: waE2E.ExtendedTextMessage_DEFAULT.Enum(),
			},
		})
		fmt.Println("err : ", err)
		fmt.Println("oke : ", oke)
		return
		// // fmt.Println(v.Message.GetConversation())
		// // client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{ButtonsMessage: &waE2E.ButtonsMessage{ContentText: proto.String("coba"), Buttons: []*waE2E.ButtonsMessage_Button{*waE2E.ButtonsMessage_Button{ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("coba2  22")}}}}})
		// // client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{

		// // 	ButtonsMessage: &waE2E.ButtonsMessage{
		// // 		ContentText: proto.String("Content"),
		// // 		FooterText:  proto.String("Footer"),
		// // 		ContextInfo: &waE2E.ContextInfo{},
		// // 		Buttons: []*waE2E.ButtonsMessage_Button{
		// // 			{
		// // 				ButtonID:       proto.String("ButtonId"),
		// // 				ButtonText:     &waE2E.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("Ok")},
		// // 				Type:           waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
		// // 				NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{},
		// // 			},
		// // 		},
		// // 	},
		// // })
		// // client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{Conversation: proto.String("cobaa")})
		// // return
		// if v.Message.GetConversation() == "#getJID" {

		// 	client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{Conversation: proto.String(v.Info.Chat.User + "@" + v.Info.Chat.Server)})
		// 	return

		// }
		// if v.Info.IsFromMe {
		// 	return
		// }
		// isipesan1 := ""
		// if v.Message.GetExtendedTextMessage() != nil {
		// 	isipesan1 = *v.Message.GetExtendedTextMessage().Text

		// } else if v.Message.GetImageMessage() != nil {

		// 	isipesan1 = v.Message.ImageMessage.GetCaption()
		// } else {
		// 	isipesan1 = v.Message.GetConversation()

		// }

		// // fmt.Println("==========================================")
		// // fmt.Println(v.Message.DocumentMessage)
		// // fmt.Println("==========================================")
		// // document := v.Message.DocumentMessage

		// // client.Download(document)
		// fmt.Println(isipesan1)
		// var nowa string
		// // nowa := "082 231 374 867"
		// if v.Info.Chat.User == "6282266353193" {
		// 	nowa = ConverNoHp("085237342776")
		// } else {
		// 	nowa = ConverNoHp(v.Info.Chat.User)
		// }
		// if strings.EqualFold(isipesan1, "order") {

		// 	pesanjawab(1, nowa, "", v.Info.Chat, nil, 0)

		// 	return
		// }
		// if readNwrite(false, nowa, 0) == 1 {
		// 	pesanjawab(2, nowa, isipesan1, v.Info.Chat, nil, 0)
		// 	// readNwrite(true, nowa, 2)
		// 	return
		// } else if readNwrite(false, nowa, 0) == 2 {
		// 	if v.Message.DocumentMessage != nil {
		// 		pesanjawab(3, nowa, isipesan1, v.Info.Chat, v.Message.DocumentMessage, 0)
		// 		// readNwrite(true, nowa, 0)
		// 		return
		// 	}

		// }
		// fmt.Println("readNwrite(false, helper.ConverNoHp(v.Info.Chat.User), 0) : ", ConverNoHp(nowa), "  :", readNwrite(false, ConverNoHp(nowa), 0))
		// if readNwrite(false, ConverNoHp(nowa), 0) == -1 {
		// 	client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{Conversation: proto.String("Ketik *order* untuk memulai transaksi")})
		// 	return
		// }
		// return
		// // fmt.Println("==============================-=-=-=-=-=-=-=--======================")
		// // document := v.Message.GetDocumentMessage()
		// // // if v.Message.DocumentMessage != nil {

		// // if document != nil {
		// // 	data, err := client.Download(document)
		// // 	if err != nil {
		// // 		fmt.Printf("Failed to download audio: %v", err)
		// // 		return
		// // 	}

		// // 	exts, _ := mime.ExtensionsByType(document.GetMimetype())
		// // 	fmt.Println(exts[0])

		// // 	if exts[0] == "xlsx" || exts[0] == "xlsm" || exts[0] == "xlsb" || exts[0] == "xltx" || exts[0] == "csv" {
		// // 		// path := fmt.Sprintf("./Downloads/Documents/%s-%s%s", v.Info.PushName, v.Info.ID, exts[0])
		// // 		path := fmt.Sprintf("./Downloads/Documents/%s", *v.Message.DocumentMessage.FileName)

		// // 		err = os.WriteFile(path, data, 0600)
		// // 		if err != nil {
		// // 			log.Printf("Failed to save document: %v", err)
		// // 			return
		// // 		}
		// // 		log.Printf("Saved document in message to %s", path)

		// // 		// 	errrror := UploadFileTest("http://localhost:8080/file", "file", path, *v.Message.DocumentMessage.FileName)
		// // 		// 	if errrror != nil {

		// // 		// 		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// // 		// 		fmt.Println(errrror)
		// // 		// 		fmt.Println("erooororoororoeroreoieroeroieroieroierio")
		// // 		// 		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// // 	}
		// // 	// 	fmt.Println("jalnanananananananan")
		// // 	// }
		// // 	return

		// // }
		// // fmt.Println("==============================-=-=-=-=-=-=-=--======================")

		// // fmt.Println("==============================-=-=-=-=-=-=-=--=")

		// // perbedaan := time.Now()
		// // difference := perbedaan.Sub(v.Info.Timestamp)

		// // if difference.Hours() >= 1 {
		// // 	fmt.Println("waktu sudah habnis")
		// // 	return
		// // }

	}

}
