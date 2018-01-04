package main

import (
	//	"bytes"
	//	"github.com/gorilla/websocket"
	"log"
	//	"net/http"
	"strings"
	// "time"
	"os"
)

var rootdir = "./"

func checkgroupadmin(groupfn string, user string) bool {

	list := readintoslice(groupfn)

	for _, li := range list {
		usern := strings.Split(li, ",")
		if user == usern[0] && usern[1] == "admin" {
			return true
		}
	}
	return false
}

func checkuseringroup(groupfn string, user string) bool {

	list := readintoslice(groupfn)

	for _, aa := range list {
		usern := strings.Split(aa, ",")
		if user == usern[0] {
			return true
		}
	}
	return false
}

func cachemsg(msg []byte) {
	dst := getjson(msg, "to")
	if dst == "group" {

		return
	}
	createuser(getjson(msg, "to"))
	uuid := getjson(msg, "uuid")

	// log.Println("caching uuid", uuid[3:],  )

	fn := rootdir + "users/" + dst + "/chatcache/msg" + uuid[3:] + ".txt"
	log.Println("caching", fn, string(msg))

	dat := []string{string(msg)}
	writeslicetofile(dat, fn)
	// err := ioutil.WriteFile(fn, msg, 0644)
	// check(err)
}

func createuser(user string) {
	mkdir(rootdir + "users/" + user + "/chatcache/")
	mkdir(rootdir + "users/" + user + "/archive/")
}

func logmsg(d1 []byte) {
	fn := rootdir + "users/log/msg" + thuman() + ".txt"
	dat := []string{string(d1)}
	writeslicetofile(dat, fn)

	// err := ioutil.WriteFile(fn, d1, 0644)
	// check(err)
}

//
//
//========================================================================
//
//
// group handlers
//
//
//========================================================================
//
//

func getuseringroup(groupfn string, user string) []string {

	list := readintoslice(groupfn)

	for _, li := range list {
		usern := strings.Split(li, ",")
		if user == usern[0] {
			return usern
		}
	}
	return nil
}

func makegroupadmin(d1 []byte) {
	from := getjson(d1, "from")
	groupname := getjson(d1, "groupname")
	adminmember := getjson(d1, "membername")
	requuid := getjson(d1, "uuid")

	groupdir := rootdir + "users/groups/" + groupname

	groupfn := groupdir + "/current.txt"
	backfn := groupdir + "/admin_" + thuman() + ".txt"

	if checkgroupadmin(groupfn, from) != true {
		clog("makegroupadmin", adminmember, groupname, from, " not admin ")
		return
	}
	if len(adminmember) < 1 {
		clog("makegroupadmin", adminmember, " not found ")
		return
	}
	// groupname := cmd[1]
	// addmember := cmd[2]
	// adminormember := cmd[3]
	// groupfn := rootdir + "users/groups/" + groupname + ".txt"

	list := readintoslice(groupfn)
	writeslicetofile(list, groupdir+"/current_"+thuman()+".txt")
	an := -1
	for li, aa := range list {
		usern := strings.Split(aa, ",")
		log.Println(li, ",", aa, ",", adminmember, ",", usern[0])
		if adminmember == usern[0] {
			log.Println(li, aa, "found")
			an = li
		}

	}
	log.Println(an)
	if an > 0 {
		list = append(list[:an], list[an+1:]...)
		list = append(list, adminmember+",admin")
		writeslicetofile(list, groupfn)
		// append2file(backfn, adminmember+",admin,"+requuid)
		lists := []string{adminmember + ",admin," + requuid}
		writeslicetofile(lists, backfn)

	}

}

func getuserinfo(d1 []byte) {

}

func getgrouplist(groupname string, from string) []string {
	// from := getjson(d1, "from")
	// groupname := group
	groupdir := rootdir + "users/groups/" + groupname
	groupfn := groupdir + "/current.txt"
	if checkuseringroup(groupfn, from) != true {

		return nil
	}

	return readintoslice(groupfn)
}

func getgrouprecipients(groupname string, from string) []string {
	// from := getjson(d1, "from")
	// groupname := group
	groupdir := rootdir + "users/groups/" + groupname
	groupfn := groupdir + "/current.txt"
	if checkuseringroup(groupfn, from) != true {
		clog("nouser")
		return nil
	}
	lines := readintoslice(groupfn)
	if len(lines) < 1 {
		clog("nolines")
		return nil
	}
	var outlist []string
	for _, li := range lines {
		usern := strings.Split(li, ",")
		// clog("append", usern[0])
		outlist = append(outlist, usern[0])
	}
	return outlist

}

func deletefromgroup(d1 []byte) {
	from := getjson(d1, "from")
	groupname := getjson(d1, "groupname")
	delmember := getjson(d1, "membername")
	adminormember := getjson(d1, "membertype")
	requuid := getjson(d1, "uuid")
	groupdir := rootdir + "users/groups/" + groupname

	groupfn := groupdir + "/current.txt"
	backfn := groupdir + "/del_" + thuman() + ".txt"

	if checkgroupadmin(groupfn, from) != true {
		clog("deletefromgroup", delmember, groupname, from, " not admin ")

		return
	}
	list := readintoslice(groupfn)
	writeslicetofile(list, groupdir+"/current_"+thuman()+".txt")
	an := -1
	for li, aa := range list {
		usern := strings.Split(aa, ",")
		log.Println(li, aa, delmember, usern[0])
		if delmember == usern[0] {
			log.Println(li, aa)
			an = li
		}

	}
	log.Println(an)
	if an > 0 {
		list = append(list[:an], list[an+1:]...)
		writeslicetofile(list, groupfn)
		smlist := []string{delmember + "," + adminormember + "," + requuid}
		writeslicetofile(smlist, backfn)
		// append2file(backfn, delmember+","+adminormember+","+requuid+"\n")

	}

}

func addtogroup(d1 []byte) {
	// log.Printf("addtogroup", d1)
	// { "type":"system", "from":"a@a.com","to":"system","msg":"addtogroup" ,"cmd": "addtogroup", "groupname":"a@a.com#grouptest", "membernae":"b@a.com", "membertype":"member" ,"time":1496789267895   ,"uuid":"msg1496789267895.YkBhLmNvbSxzeXN0ZW0sYWRkdG9ncm91cA=="})
	from := getjson(d1, "from")
	groupname := getjson(d1, "groupname")
	addmember := getjson(d1, "membername")
	adminormember := getjson(d1, "membertype")
	requuid := getjson(d1, "uuid")
	groupdir := rootdir + "users/groups/" + groupname
	groupfn := groupdir + "/current.txt"
	backfn := groupdir + "/add_" + thuman() + ".txt"

	if checkgroupadmin(groupfn, from) != true {
		clog("addtogroup", addmember, groupname, from, " not admin ")
		return
	}

	if len(getuseringroup(groupfn, addmember)) > 1 {
		clog("addtogroup", groupfn, addmember, "already in")
		//	fwdmessage([]byte(`{"to": "` + from + `","msg":"already in", "type":"system", "uuid": "addtogroupfail"  ,"replyto":"` + requuid + `"}`))
		return
	}

	if _, err := os.Stat(groupfn); err == nil {
		// append2file(groupfn, addmember+","+adminormember+","+requuid+"\n")
		// append2file(backfn, addmember+","+adminormember+","+requuid+"\n")
		list := readintoslice(groupfn)
		list = append(list, addmember+",member")
		writeslicetofile(list, groupfn)
		smlist := []string{addmember + "," + adminormember + "," + requuid}
		writeslicetofile(smlist, backfn)
		//		fwdmessage([]byte(`{"to": "` + from + `","msg":"addtogroup", "type":"system", "uuid": "addtogroupok"  ,"replyto":"` + requuid + `"}`))
	} else {
		log.Println("addtogroup", groupfn, "no file")
	}
}

func creategroup(d1 []byte) {

	from := getjson(d1, "from")
	groupname := getjson(d1, "groupname")
	requuid := getjson(d1, "uuid")
	if len(from) == 0 {
		return
	}
	if strings.Split(groupname, "#")[0] != from {
		return
	}

	log.Printf("creategroup", requuid)

	groupdir := rootdir + "users/groups/" + groupname
	mkdir(groupdir)
	groupfn := groupdir + "/current.txt"
	backfn := groupdir + "/current_" + thuman() + ".txt"

	if _, err := os.Stat(groupfn); err == nil {
		clog("creategroup", groupfn, "exists")
		//	fwdmessage([]byte(`{"from":"system","to": "` + from + `","msg":"creategroup,exists,` + groupname + `", "type":"system", "uuid": "creategroupexists" ,"replyto":"` + requuid + `"}`))
	} else {
		clog("creategroup", groupfn, "")

		list := []string{from + ",admin"}
		writeslicetofile(list, groupfn)
		writeslicetofile(list, backfn)
		// append2file(groupfn, from+",admin\n")
		// append2file(backfn, from+",admin\n")
		//fwdmessage([]byte(`{"from":"system","to": "` + from + `","msg":"creategroup,ok,` + groupname + `", "type":"system", "uuid": "creategroupok","replyto":"` + requuid + `"}`))
	}

}
