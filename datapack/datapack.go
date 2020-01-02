package datapack

import (
	"fmt"
	"io/ioutil"
	"net/http"
    "github.com/likexian/whois-go"
    "encoding/json"//
	"bufio"
	s "strings"
    "bytes"
	"regexp"
	"time"
	"database/sql"
    _ "github.com/lib/pq"
    "sort"
)

type Item struct{
  Domain string `json:"domain"`
  Info Server_data `json:"info"`
}

type Items struct{
  Items_data []Item `json:"items"`
}

type endpoint struct{
  IpAddress string `json:"ipAddress"`
  ServerName string `json:"serverName"`
  Grade string  `json:"grade"`
  OrgName string `json:"orgName"`
  Country string `json:"country"`
}

type Server_data struct{
  Endpoint []endpoint `json:"endpoints"`
  ServersChanged bool `json:"servers_changed"`
  SslGrade string `json:"ssl_grade"`
  PreviousSslGrade string `json:"previous_ssl_grade"`
  Logo string `json:"logo"`
  Title string `json:"title"`
  IsDown bool `json:"is_down"`		   
}





func Data_obtain(domain string, url string, server_data *Server_data){
	t := time.Now().Minute()
	resp, err := http.Get(url)
	if err != nil {
		server_data.IsDown = true
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile("test.json", body, 0644)
        json.Unmarshal(body, &server_data)
	domain_data_obtainer(server_data)

	if resp.StatusCode == 200 {
		server_data.IsDown = false
		html_data_obtainer(string("https://"+domain), server_data)	
	}else{server_data.IsDown = true}

	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/bank?sslmode=disable")
        if err != nil {
          panic(err)
        }

	 _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS server_db (id serial PRIMARY KEY, domain VARCHAR(20),time INT,info json);
     `)
       if  err != nil {
          panic(err)
       }

       b,_ := json.Marshal(server_data)
   
       _, err = db.Exec(`
           INSERT INTO server_db (domain,time,info) VALUES($1,$2,$3)`,domain,t,b)
      if err != nil {
         panic(err)
      }

    com_db(db, domain, server_data)    
	saver_db(db,domain,server_data) 	
}

func com_db(db *sql.DB, domain string,server_data *Server_data){
	
	var newAfter string
    var After string

	rows, err := db.Query(`
    SELECT info FROM server_db 
    WHERE domain = $1 AND time = (SELECT MIN(time) FROM server_db)`, domain)
    if err != nil {
       panic(err)
    }
    defer rows.Close()
	if rows.Next(){
    rows.Scan(&newAfter)       
    } 

    rows, err = db.Query(`
    SELECT info FROM server_db 
    WHERE domain = $1 AND time = (SELECT MAX(time) FROM server_db)`, domain)
    if err != nil {
       panic(err)
    }
    defer rows.Close()
	if rows.Next(){
    rows.Scan(&After)       
    } 
    if rows.Next(){
    rows.Scan(&After)       
    }

    server_data.ServersChanged = After==newAfter

    var before, after Server_data      
    json.Unmarshal([]byte(newAfter), &before)
    json.Unmarshal([]byte(After), &after)
    
    if(after.Endpoint != nil){
    	server_data.SslGrade = selector_grade(&after)
    }

    if(before.Endpoint != nil){
    	server_data.PreviousSslGrade = selector_grade(&before)
    }

    
}

func saver_db(db *sql.DB, domain string,server_data *Server_data){

	 _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS q_db (id serial PRIMARY KEY, domain VARCHAR(20),info json);
     `)
       if  err != nil {
          panic(err)
       }
    
       b,_ := json.Marshal(server_data)
   
       _, err = db.Exec(`
           INSERT INTO q_db (domain,info) VALUES($1,$2)`,domain,b)
      if err != nil {
         panic(err)
      }
}

func domain_data_obtainer (server_data *Server_data){

	for i := 0; i < len(server_data.Endpoint); i++ {	
		out, err := whois.Whois(server_data.Endpoint[i].IpAddress)
		if err != nil {
	    		fmt.Println(err)
		}
		read := bytes.NewReader([]byte(out))
		scanner := bufio.NewScanner(read)
		var j int
		even := 0
		for scanner.Scan(){
   			j++ 
			linea := scanner.Text()
			
			if s.Contains(linea,"OrgName"){
				elem := s.Index(linea,":")+1
				OrgName := s.Replace(string(linea[elem:]),"  ", "",444)
				server_data.Endpoint[i].OrgName = OrgName 
				even++
			}

			if s.Contains(linea,"Country"){
				elem := s.Index(linea,":")+1
				Country := s.Replace(string(linea[elem:]),"  ", "",444)
				server_data.Endpoint[i].Country = Country
				even++
			}
		if even == 2{break}
		}	
	}



}

func html_data_obtainer(url string, server_data *Server_data){

	resp, err := http.Get(url)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()



	html, _ := ioutil.ReadAll(resp.Body)
        read := bytes.NewReader(html)
	scanner := bufio.NewScanner(read)

	var j int
	for scanner.Scan(){
		j++ 
		linea := scanner.Text()

		if s.Contains(linea,"<meta charset=\"utf-8\"/>"){    
			re := regexp.MustCompile(`<title.*?>(.*)</title>`)
			submatchall := re.FindAllStringSubmatch(linea, -1)
			server_data.Title = submatchall[len(submatchall)-1][1]
			re = regexp.MustCompile(`<link href=[\"'](.+?)[\"'].*?>`)
			submatchall = re.FindAllStringSubmatch(linea, -1)
			server_data.Logo = submatchall[len(submatchall)-1][1]
			break
		}
	}	
}

func selector_grade(server_data *Server_data)string{

	var st []string
	for _, elm := range server_data.Endpoint{
		if(len(elm.Grade) == 0){st = append(st,"_")
		}else{st = append(st,elm.Grade)}
	}
	sort.Strings(st)
	r := st[0] 
	for _, elm := range st{
		if s.Contains(elm,r) && s.Contains(elm,r+"+"){
			return elm
		}
	}
	
	return r

}

func Builder_hist_db(items_data *Items){

	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/bank?sslmode=disable")
        if err != nil {
          panic(err)
        }

rows, err := db.Query(`
    SELECT domain,info FROM q_db 
    `)
    if err != nil {
       panic(err)
    }
    defer rows.Close()
var q_down []Server_data
var domains_save []string
i := 0
	for rows.Next(){
		var queries string
                var domain string
		var sd Server_data
		rows.Scan(&domain,&queries)
		if err != nil{
			panic(err)
		}

	json.Unmarshal([]byte(queries), &sd)
	q_down = append(q_down,sd)
        domains_save = append(domains_save,domain)	
	
	i++
	}	



for i:=0;i<len(q_down);i++{

var i_data Item
i_data.Domain = domains_save[i]
i_data.Info = q_down[i]

items_data.Items_data = append(items_data.Items_data, i_data)

}

 _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS hist_db (id serial PRIMARY KEY, info json);
     `)
       if  err != nil {
          panic(err)
       }
    
       b,_ := json.Marshal(items_data)
   
       _, err = db.Exec(`
           INSERT INTO hist_db (info) VALUES($1)`,b)
      if err != nil {
         panic(err)
      }

}