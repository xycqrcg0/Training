from bs4 import BeautifulSoup, Comment
import requests
import html2text
import re
import os
import time

def get_post2markdown(url,t):
    print(url)
    headers = {
        "Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
        "User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
        "Cache-Control":"max-age=0"
    }
    response = requests.get(url,headers=headers)
    response.encoding = 'utf-8'
    soup = BeautifulSoup(response.text,"html.parser")
    ht = html2text.HTML2Text()
    ht.body_width = 0
    
    # 首先，文章内容在article标签下 
    article = soup.find("article")
    # 处理一下代码块
    for figure in article.find_all("figure",class_="highlight") :
        table = figure.find("table")
        if not table :
            continue
        
        new_pre = soup.new_tag("pre")
        new_code = soup.new_tag("code")
        new_pre.append(new_code)
        
        code_td = table.find("td",class_="code")
        if not code_td :
            continue
        
        if len(figure) > 1 :    
            lang = figure["class"][1]
            new_code.append("``` "+lang+"\n")
        else :
            new_code.append("```\n")


        for line in code_td.find_all("span",class_="line"):
            line_text = line.get_text().replace("\n"," ")
            new_code.append(line_text + "\n")
    
        new_code.append("``` \n")
    
        table.replace_with(new_pre)    
    
    # 删除结尾处的一些链接
    dele= set()
    next = article.find("hr",style="margin-top: 2rem;")
    while next :
        dele.add(next)
        next = next.next_sibling
    for tag in dele:
        tag.decompose()
    
    
    title_tag = article.find("h1")
    title = title_tag.text
    # 还要防止title里出现“/”
    if "/" in title:
        title = title.replace("/","_")

    # print(str(article))
    text = ht.handle(str(article)).strip()
    re.sub(r'\n{2,}', '\n', text.strip())
    # print(text)
    
    cur = os.getcwd()
    path = os.path.join(cur,"geektutu",t[:7])
    if not os.path.exists(path):
        os.mkdir(path)

    file_name = t[8:10]+"-"+title+".md"
    file_path = os.path.join(path,file_name)
    with open(file_path,"a+") as file :
        file.write(text)
    
    print("finish one.")


def get_link(p_url,url):
    headers = {
        "Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
        "User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
        "Cache-Control":"max-age=0"
    }
    response = requests.get(url,headers=headers)
    response.encoding = 'utf-8'
    soup = BeautifulSoup(response.text,"html.parser")
    ht = html2text.HTML2Text()
    ht.body_width = 0

    tag = soup.find("div",class_="float-left post-container box-shadow")
    # 先把置顶文章删了
    comments = soup.find_all(string = lambda string : isinstance(string,Comment))
    dele = set()
    for comment in comments :
        if comment.strip() == "置顶文章" :
            next = comment.next_sibling
            while next :
                if isinstance(next,Comment) and next.strip() == "普通文章" :
                    break
                dele.add(next)
                next = next.next_sibling
    for d in dele :
        d.decompose()
    
    posts = tag.find_all("div",class_="post-preview")
    for post in posts :
        link_tag = post.find("a",class_="title")
        link = link_tag["href"]
        post_url = p_url + "/" + link
        time_tag = post.find("small")
        time_str = time_tag.text
        t = time_str[4:14]
        get_post2markdown(post_url,t)
        time.sleep(0.5)


p_url = "https://geektutu.com"
common_url = "https://geektutu.com/page"


def main ():
    cur = os.getcwd()
    path = os.path.join(cur,"geektutu")
    if not os.path.exists(path):
        os.mkdir(path)

    # 先爬第一页
    get_link(p_url,p_url)

    pagei = 2
    while True:
        page = str(pagei)
        url = common_url+"/"+page+"/"

        get_link(p_url,url)

        pagei += 1
        if pagei >= 12:
            break
        
        time.sleep(0.5)

if __name__ == "__main__":
    main()