package main

import (
	"fmt"
//	"bytes"
	"io/ioutil"
	"net/http"
        "github.com/likexian/whois-go"
        "os"
        "encoding/json"//
	"bufio"
	s "strings"
        //"github.com/likexian/whois-parser-go"
//    "github.com/go-chi/chi/middleware"
//    _ "github.com/go-sql-driver/mysql"
)

const(
domain = "truora.com"
)

type endpoint struct{
  IpAddress string `json:"ipAddress"`
  ServerName string `json:"serverName"`
  StatusMessage string `json:"statusMessage"`
  StatusDetails string `json:"statusDetails"`
  StatusDetailsMessage string `json:"statusDetailsMessage"`
  Grade string  `json:"grade"`
  GradeTrustIgnored string `json:"gradeTrustIgnored"`
  HasWarnings bool `json:"hasWarnings"`
  IsExceptional bool `json:"isExceptional"`
  Progress int `json:"progress"`
  Eta int `json:"eta"`
  Duration int  `json:"duration"`
  Delegation int `json:"delegation"`
  OrgName string `json:"orgName"`
  Country string `json:"country"`
}

type Server_data struct{
  Host string `json:"host"`
  Port int `json:"port"`
  Protocol string `json:"protocol"`
  IsPublic bool `json:"isPublic"`
  Status string `json:"status"`
  StartTime int `json:"startTime"`
  EngineVersion string `json:"engineVersion"`
  CriteriaVersion string `json:"criteriaVersion"`
  Endpoint []endpoint `json:"endpoints"`
  Logo string `json:"logo"`
  Title string `json:"title"`	   
}


func main() {
	url := "https://api.ssllabs.com/api/v3/analyze?host=" + domain
        var server_data Server_data
	data_obtain(url, &server_data)
	domain_data_obtainer(url, &server_data)
	html_data_obtainer(string("https://"+domain), &server_data)
	
	fmt.Println(server_data)
	

}

func data_obtain(url string, server_data *Server_data){
	
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		fmt.Println("Error in importing")
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("test.json", body, 0644)
        json.Unmarshal(body, &server_data)

 
}

func domain_data_obtainer (url string, server_data *Server_data){

	for i := 0; i < len(server_data.Endpoint); i++ {
    		
		out, err := whois.Whois(server_data.Endpoint[i].IpAddress)
		if err != nil {
	    		fmt.Println(err)
		}
		
		f,_ := os.Create("out.text")
		f.WriteString(out)
		defer f.Close()

		archivo,_ := os.Open("./out.text")
		scanner := bufio.NewScanner(archivo)
		defer archivo.Close()		

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
	fmt.Println("Error noo!!")
	panic(err)
	}

	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	
	ioutil.WriteFile("page.html", html, 444)
	
	archivo,_ := os.Open("./page.html")
	scanner := bufio.NewScanner(archivo)
	defer archivo.Close()
	
	var j int
	for scanner.Scan(){
		j++ 
		linea := scanner.Text()

		if s.Contains(linea,"<!DOCTYPE html>"){    		
			elem := s.Index(linea,"<title>")+len("<title>")
			elem1 := s.Index(linea,"</title>")			
			server_data.Title = string(linea[elem:elem1])
			elem = s.Index(linea,"</script><link href=")+len("</script><link href=")
			elem1 = s.Index(linea," rel=\"shortcut icon\" type=\"image/x-icon\"/>")
			server_data.Logo = string(linea[elem:elem1])
			break
		}
	}	
}



