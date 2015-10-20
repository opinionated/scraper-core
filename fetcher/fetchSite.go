package fetcher

import (
	"encoding/xml"
	"fmt"
	"github.com/opinionated/scraper-core/scraper"
	"golang.org/x/net/html"
	"net/http"
)

// Parse the wsj opinion rss feed into articles.
// Outputs articles and an error
func GetStories(rss scraper.RSS, body []byte) error {
	err := xml.Unmarshal(body, rss)
	if err != nil {
		fmt.Printf("err:", err)
		return err
	}

	for i := 0;
	 i < rss.GetChannel().GetNumArticles();
	 i++ {
		article := rss.GetChannel().GetArticle(i)
		fmt.Println("index:", i, "title:", article.GetTitle())
		 // "\tdescr:", article.GetDescription())
	}

	return nil
}

// Request a page containing the article linked to
func DoGetArticle(article scraper.Article) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", article.GetLink(), nil)			//create http request
	req.Header.Add("Referer", "https://www.google.com")					//required to get past paywall

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8") //extra?
	req.Header.Add("Accept-Language", "en-US,en;q=0.5") 				//extra?
	req.Header.Add("Host", "www.wsj.com") 								//extra?
	req.Header.Add("Cookie", "DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dalbany") 								//messin with cookies

	resp, err := client.Do(req)											//send http request
	if err != nil {
		fmt.Println("oh nose, err with get article http request:", err)
		return err
	}
	if (resp.Header["X-Article-Template"][0] != "full"){				//error checking 
		fmt.Println("OH GOD OH GOD OH GOD, the template isn't full THEY KNOW\n", err)	//Should panic if this isn't at full
		return err
	}	
	fmt.Println("---------------------------------RESPONSE\n")			//check to see if X-Article-Template is [full]
	fmt.Println(resp.Header)
	fmt.Println("---------------------------------RESPONSE\n")

// DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dalbany


	defer resp.Body.Close()

	parser := html.NewTokenizer(resp.Body)
	tmp := article.(*scraper.WSJArticle)
	tmp.DoParse(parser)													//parse the html body
	return err
}
// ToDo, make some cookie profiles and make a function return a random cookie string

	//chrome on windows 8.1
// dji_user=cd1f7b93-5270-4177-a381-864230207766; fyre-livecount=687388676917; DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dtroy%7C%7Clatitude%3D42.7478%7C%7Clongitude%3D-73.6049%7C%7Ctimezone%3Dest%7C%7Czip%3D12180-12183%7C%7CORCS%3Dna%2Cus; test_key=0.7395806894637644; D_SID=129.161.129.244:UnlejPUojYIvjNrUdoiT2cD2zcZ8GNIkA1tKARD+BA4; _cb_ls=1; __qca=P0-619642729-1444422342960; cX_S=ifk3rot7xpieb2nq; s_vnum=1475958342381%26vn%3D2; s_vmonthnum=1446350400385%26vn%3D2; dji_user=d54eda3b-a93a-4763-8c1f-304929d88909; DJCOOKIE=ORC%3Dna%2Cus%7C%7CFCFGOOGLE%3D0%7C%7CFCFEXPGOOGLE%3D1444674503096; wsjregion=na%2Cus; D_PID=469F0452-18FF-3E05-8072-566D9785BE96; D_IID=007E6ABE-A48C-3DE5-81E0-CBECBC2C96AB; D_UID=51DCB921-500A-36DD-9325-5584F17E13D0; D_HID=Jxzp7NvLL6qaRXhQ+xTYS1cJRFKZBRcxjajzb0KR7nA; utag_main=v_id:01504e47d6880002494f77ff63b40606d005a065009dc$_sn:2$_ss:0$_st:1444589906195$_pn:2%3Bexp-session$ses_id:1444588095632%3Bexp-session; s_cc=true; s_fid=429C8463B968A52F-2BB5BC9D3004ED8F; s_sq=%5B%5BB%5D%5D; _chartbeat2=DsofJnZWBgXNDUuc.1441222350395.1444588107532.0000000000000101; cX_P=ie36kra7n61qjwjm; bkuuid=LQ80JkoD99O5TgjP
// DJCOOKIE=ORC%3Dna%2Cus%7C%7CFCFGOOGLE%3D0%7C%7CFCFEXPGOOGLE%3D1444674503096;
// D_SID=129.161.129.244:UnlejPUojYIvjNrUdoiT2cD2zcZ8GNIkA1tKARD+BA4;
// D_PID=469F0452-18FF-3E05-8072-566D9785BE96;
// D_IID=007E6ABE-A48C-3DE5-81E0-CBECBC2C96AB;
// D_UID=51DCB921-500A-36DD-9325-5584F17E13D0;
// D_HID=Jxzp7NvLL6qaRXhQ+xTYS1cJRFKZBRcxjajzb0KR7nA;
// test_key=0.7395806894637644;
// s_vnum=1475958342381%26vn%3D2;
// s_vmonthnum=1446350400385%26vn%3D2;
// _chartbeat2=DsofJnZWBgXNDUuc.1441222350395.1444588107532.0000000000000101;
// s_fid=429C8463B968A52F-2BB5BC9D3004ED8F;
// bkuuid=LQ80JkoD99O5TgjP
// cX_P=ie36kra7n61qjwjm;
// __qca=P0-619642729-1444422342960;
// utag_main=v_id:01504e47d6880002494f77ff63b40606d005a065009dc$_sn:2$_ss:0$_st:1444589906195$_pn:2%3Bexp-session$ses_id:1444588095632%3Bexp-session;
//!// wsjregion=na%2Cus;
//!// _cb_ls=1;

// dji_user=cd1f7b93-5270-4177-a381-864230207766;
// fyre-livecount=687388676917;
// DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dtroy%7C%7Clatitude%3D42.7478%7C%7Clongitude%3D-73.6049%7C%7Ctimezone%3Dest%7C%7Czip%3D12180-12183%7C%7CORCS%3Dna%2Cus;
// cX_S=ifk3rot7xpieb2nq;
// dji_user=d54eda3b-a93a-4763-8c1f-304929d88909;
// s_cc=true;
// s_sq=%5B%5BB%5D%5D;

	//firefox on ubuntu linux
// djcs_route=dc9531e6-4e8b-42e1-9697-d32127121e7a; mmcore.tst=0.143; mmid=1442230749%7CAQAAAAoB73ndgwwAAA%3D%3D; mmcore.pd=601088381%7CAQAAAAoBQgHved2DDGFQrF4BAM7Xulbd0NJIDwAAAM7Xulbd0NJIAAAAAAEAAAD/////AA53d3cuZ29vZ2xlLmNvbQSDDAEAAAAAAAAAAQAA////////////////AQAiRwAAAC6bh2GDDAD/////AYMMgwz//wEAAAEAAAAAAROvAABuFQEAAAAAAAFF; mmcore.srv=nycvwcgus10; mm_criteria=%7B%22ReferringSource%22%3A%22Google%20Web%22%2C%22OnsiteChannel%22%3A%22%22%7D; utag_main=v_id:01504e04ed6100205a485f7c43ec0504c00ac009009dc$_sn:5$_ss:0$_st:1444960693124$_pn:3%3Bexp-session$ses_id:1444957583800%3Bexp-session; s_vnum=1475953957319%26vn%3D5; s_vmonthnum=1446350400322%26vn%3D5; s_fid=0E3724C4B678E00E-2E630810713E8479; cX_P=ifk15p1ljmu94v39; DJCOOKIE=ORC%3Dna%2Cus%7C%7CFCFWSJ%3D0%7C%7CFCFEXPWSJ%3D1445043979059%7C%7CFCFGOOGLE%3D1%7C%7CFCFEXPGOOGLE%3D1445044395238; wsjregion=na%2Cus; __qca=P0-1878947270-1444418008205; bkuuid=u8fdEGLr99YY%2F4hP; test_key=0.33328365022316575; D_SID=129.161.196.174:st+Zf3/otXrEk3j2TldakPKGiBwCrearhc3redl7YOU; D_PID=7E87B955-84EB-3578-A991-B8948732DC33; D_IID=D64EDFDB-4945-3C82-9C5F-E399081516C9; D_UID=13EAE792-5FBD-3AA6-857B-71A17D2FC8CD; D_HID=3tczQ2cH5YWWNyKmNf1oZQ+gU36z/P3V0XBMxgM1PpQ; __gads=ID=5b779616d5fd18ee:T=1444418012:S=ALNI_MarRU2ugjCL_8Ac8mdXjFkKRX8wIw; _cb_ls=1; _chartbeat2=PFGgwEn4U4B-VCaN.1444418016635.1444958895887.1010101; cke=%7B%22A%22%3A11%7D; s_invisit=true; s_monthinvisit=true; gpv_pn=WSJ_Opinion_Article_20151014_Democrats%20Say%20the%20Economy%20Stinks
// DJCOOKIE=ORC%3Dna%2Cus%7C%7CFCFWSJ%3D0%7C%7CFCFEXPWSJ%3D1445043979059%7C%7CFCFGOOGLE%3D1%7C%7CFCFEXPGOOGLE%3D1445044395238;
// D_SID=129.161.196.174:st+Zf3/otXrEk3j2TldakPKGiBwCrearhc3redl7YOU;
// D_PID=7E87B955-84EB-3578-A991-B8948732DC33;
// D_IID=D64EDFDB-4945-3C82-9C5F-E399081516C9;
// D_UID=13EAE792-5FBD-3AA6-857B-71A17D2FC8CD;
// D_HID=3tczQ2cH5YWWNyKmNf1oZQ+gU36z/P3V0XBMxgM1PpQ;
// test_key=0.33328365022316575;
// s_vnum=1475953957319%26vn%3D5;
// s_vmonthnum=1446350400322%26vn%3D5;
// _chartbeat2=PFGgwEn4U4B-VCaN.1444418016635.1444958895887.1010101;
// s_fid=0E3724C4B678E00E-2E630810713E8479;
// bkuuid=u8fdEGLr99YY%2F4hP;
// cX_P=ifk15p1ljmu94v39;
// __qca=P0-1878947270-1444418008205;
// utag_main=v_id:01504e04ed6100205a485f7c43ec0504c00ac009009dc$_sn:5$_ss:0$_st:1444960693124$_pn:3%3Bexp-session$ses_id:1444957583800%3Bexp-session;
//!// wsjregion=na%2Cus;
//!// _cb_ls=1;

// djcs_route=dc9531e6-4e8b-42e1-9697-d32127121e7a;
// mmcore.tst=0.143;
// mmid=1442230749%7CAQAAAAoB73ndgwwAAA%3D%3D;
// mmcore.pd=601088381%7CAQAAAAoBQgHved2DDGFQrF4BAM7Xulbd0NJIDwAAAM7Xulbd0NJIAAAAAAEAAAD/////AA53d3cuZ29vZ2xlLmNvbQSDDAEAAAAAAAAAAQAA////////////////AQAiRwAAAC6bh2GDDAD/////AYMMgwz//wEAAAEAAAAAAROvAABuFQEAAAAAAAFF;
// mmcore.srv=nycvwcgus10;
// mm_criteria=%7B%22ReferringSource%22%3A%22Google%20Web%22%2C%22OnsiteChannel%22%3A%22%22%7D;
// __gads=ID=5b779616d5fd18ee:T=1444418012:S=ALNI_MarRU2ugjCL_8Ac8mdXjFkKRX8wIw;
// cke=%7B%22A%22%3A11%7D;
// s_invisit=true;
// s_monthinvisit=true;
// gpv_pn=WSJ_Opinion_Article_20151014_Democrats%20Say%20the%20Economy%20Stinks


//response set cookie
// "DJCOOKIE=ORC%3Dna%2Cus%7C%7CFCFWSJ%3D0%7C%7CFCFEXPWSJ%3D1445043979059%7C%7CFCFGOOGLE%3D2%7C%7CFCFEXPGOOGLE%3D1445044395238;
// Domain=.wsj.com;
// Path=/;
// Expires=Sat, 15 Oct 2016 01:50:16 GMTwsjregion=na%2Cus;
// Domain=.wsj.com;
// Path=/;
// Expires=Sun, 15 Nov 2015 01:50:16 GMTDJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dtroy%7C%7Clatitude%3D42.7478%7C%7Clongitude%3D-73.6049%7C%7Ctimezone%3Dest%7C%7Czip%3D12180-12183;
// Domain=.wsj.com;
// Path=/"