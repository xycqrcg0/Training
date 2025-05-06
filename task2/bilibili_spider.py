from bs4 import BeautifulSoup
from functools import reduce
from hashlib import md5
import urllib.parse
import time
import requests
import json
import random

"""ps: 前面两个算法粘贴自 [https://github.com/SocialSisterYi/bilibili-API-collect]"""

"""bv号与av号的转换"""
XOR_CODE = 23442827791579
MASK_CODE = 2251799813685247
MAX_AID = 1 << 51
ALPHABET = "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf"
ENCODE_MAP = 8, 7, 0, 5, 1, 3, 2, 4, 6
DECODE_MAP = tuple(reversed(ENCODE_MAP))

BASE = len(ALPHABET)
PREFIX = "BV1"
PREFIX_LEN = len(PREFIX)
CODE_LEN = len(ENCODE_MAP)


def av2bv(aid: int) -> str:
    bvid = [""] * 9
    tmp = (MAX_AID | aid) ^ XOR_CODE
    for i in range(CODE_LEN):
        bvid[ENCODE_MAP[i]] = ALPHABET[tmp % BASE]
        tmp //= BASE
    return PREFIX + "".join(bvid)


def bv2av(bvid: str) -> int:
    assert bvid[:3] == PREFIX

    bvid = bvid[3:]
    tmp = 0
    for i in range(CODE_LEN):
        idx = ALPHABET.index(bvid[DECODE_MAP[i]])
        tmp = tmp * BASE + idx
    return (tmp & MASK_CODE) ^ XOR_CODE


assert av2bv(111298867365120) == "BV1L9Uoa9EUx"
assert bv2av("BV1L9Uoa9EUx") == 111298867365120

"""WBI签名算法"""
mixinKeyEncTab = [
    46,
    47,
    18,
    2,
    53,
    8,
    23,
    32,
    15,
    50,
    10,
    31,
    58,
    3,
    45,
    35,
    27,
    43,
    5,
    49,
    33,
    9,
    42,
    19,
    29,
    28,
    14,
    39,
    12,
    38,
    41,
    13,
    37,
    48,
    7,
    16,
    24,
    55,
    40,
    61,
    26,
    17,
    0,
    1,
    60,
    51,
    30,
    4,
    22,
    25,
    54,
    21,
    56,
    59,
    6,
    63,
    57,
    62,
    11,
    36,
    20,
    34,
    44,
    52,
]
# ~~这格式化真是太好用辣，一定要一行一个是吧~~

def getMixinKey(orig: str):
    "对 imgKey 和 subKey 进行字符顺序打乱编码"
    return reduce(lambda s, i: s + orig[i], mixinKeyEncTab, "")[:32]


def encWbi(params: dict, img_key: str, sub_key: str):
    "为请求参数进行 wbi 签名"
    mixin_key = getMixinKey(img_key + sub_key)
    curr_time = round(time.time())
    params["wts"] = curr_time  # 添加 wts 字段
    params = dict(sorted(params.items()))  # 按照 key 重排参数
    # 过滤 value 中的 "!'()*" 字符
    params = {
        k: "".join(filter(lambda chr: chr not in "!'()*", str(v)))
        for k, v in params.items()
    }
    query = urllib.parse.urlencode(params)  # 序列化参数
    wbi_sign = md5((query + mixin_key).encode()).hexdigest()  # 计算 w_rid
    params["w_rid"] = wbi_sign
    return params


def getWbiKeys() -> tuple[str, str]:
    "获取最新的 img_key 和 sub_key"

    resp = requests.get("https://api.bilibili.com/x/web-interface/nav", headers=h())
    resp.raise_for_status()
    json_content = resp.json()
    img_url: str = json_content["data"]["wbi_img"]["img_url"]
    sub_url: str = json_content["data"]["wbi_img"]["sub_url"]
    img_key = img_url.rsplit("/", 1)[1].split(".")[0]
    sub_key = sub_url.rsplit("/", 1)[1].split(".")[0]
    return img_key, sub_key


ch = 0
eh = 0 #记录一下
headers1 = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
    "Referer": "https://www.bilibili.com/",
}
headers2 = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
    "Referer": "https://www.bilibili.com/",
}
def h() -> dict:
    global ch, eh
    a = random.randint(0,1)
    if a == 0 :
        ch += 1
        return headers1
    else :
        eh += 1
        return headers2


web_location = 1315875
re_web_location = 333788
re_url = "https://api.bilibili.com/x/v2/reply/wbi/main?"
re_re_url = "https://api.bilibili.com/x/v2/reply/reply?"
img_key, sub_key = getWbiKeys()
MINLIKE = 1000
# 设置两个默认值吧
bvid = "BV1Ut411v74a"
oid = 38096452


def main():
    # 获取视频信息
    video_url = f"https://www.bilibili.com/video/{bvid}/"
    basic_res = requests.get(video_url, headers=headers1)
    basic_res.encoding = "utf-8"
    soup = BeautifulSoup(basic_res.text, "html.parser")

    element1 = soup.find("div", id="viewbox_report")
    # title_element = element1.find("div",class_="vedio-info-title") 为什么这个就找不到？
    title_element = element1.find("h1", class_="video-title special-text-indent")
    title = str(title_element.string)
    time_element = element1.find("div", class_="pubdate-ip-text")
    lauch_time = str(time_element.string)

    element2 = soup.find("div", class_="video-toolbar-left-main")
    video_like = str(
        element2.find("span", class_="video-like-info video-toolbar-item-text").string
    )
    video_coin = str(
        element2.find("span", class_="video-coin-info video-toolbar-item-text").string
    )
    video_fav = str(
        element2.find("span", class_="video-fav-info video-toolbar-item-text").string
    )
    video_share = str(
        element2.find("span", class_="video-share-info video-toolbar-item-text").string
    )

    element3 = soup.find("div", class_="up-detail-top")
    up_name_element = element3.find("a", href=lambda x: "space.bilibili.com" in x)
    up_name = up_name_element.get_text(strip=True)  # 直接.string好像不太行？

    print("\n视频信息:")
    print(f"<<{title}>>")
    print(f"作者:<{up_name}> ,at {lauch_time}")
    print(
        f"点赞:{video_like}     投币:{video_coin}     \n收藏:{video_fav}     转发:{video_share} \n\n"
    )
    print("评论信息：")

    emojis = {}
    offset = '{"offset":""}'
    pa = 1
    while True:
        br = True
        # 啊啊每个请求的wbi签名都不一样！每次要重新生成！（被硬控了...）
        signed_params = encWbi(
            params={
                "oid": oid,
                "type": 1,
                "mode": 3,
                "pagination_str": offset,
                "plat": 1,
                "web_location": web_location,
            },
            img_key=img_key,
            sub_key=sub_key,
        )
        query = urllib.parse.urlencode(signed_params)
        response = requests.get(re_url + query, headers=h())
        response.encoding = "utf-8"
        # 获得的是json数据
        data = json.loads(response.text)

        raw_cursor = data["data"]["cursor"]["pagination_reply"]["next_offset"]
        offset = '{"offset":"'
        for char in raw_cursor:
            if char == '"':
                offset += '\\"'
            else:
                offset += char
        offset += '"}'

        replies = data["data"]["replies"]
        for reply in replies:
            like = reply["like"]
            if like >= MINLIKE:
                br = False

                ctime = reply["ctime"]
                s_ctime = time.localtime(ctime)
                t = time.strftime("%Y-%m-%d %H:%M:%S", s_ctime)
                name = reply["member"]["uname"]
                msg = reply["content"]["message"]

                print(f"## {name}:  (at {t}, 点赞数：{like})")
                print(f"   {msg}")
                emotes = reply["content"].get("emote")  # 该成员可能不存在,因此要用get()
                if emotes is not None:
                    for emote in emotes:
                        if emote[0] == "[":
                            ok = emojis.get(emote)
                            if ok is not None:
                                emojis[emote] += 1
                            else:
                                emojis[emote] = 1
            # 该评论下的子评论
            root = reply["rpid"]
            # 确定一下最多有多少页子评论(该成员可能不存在)
            rnumstr = reply["reply_control"].get("sub_reply_entry_text")
            if rnumstr == None:  # 没子评论
                continue

            rnum = int(rnumstr[1:3])
            pages = (rnum + 19) // 10
            pn = 1
            while pn <= pages:
                # 如果有一页的评论都没有1000个赞，之后的评论就不处理了~
                exit = True

                re_params = {
                    "oid": oid,
                    "type": 1,
                    "root": root,
                    "ps": 20,
                    "pn": pn,
                    "web_location": re_web_location,
                }
                query = urllib.parse.urlencode(re_params)
                res = requests.get(re_re_url + query, headers=h())
                res.encoding = "utf-8"
                re_data = json.loads(res.text)
                re_replies = re_data["data"]["replies"]
                for reply in re_replies:
                    like = reply["like"]
                    if like >= MINLIKE:
                        exit = False

                        ctime = reply["ctime"]
                        s_ctime = time.localtime(ctime)
                        t = time.strftime("%Y-%m-%d %H:%M:%S", s_ctime)
                        name = reply["member"]["uname"]
                        msg = reply["content"]["message"]

                        print(f"   ### {name}:  (at {t}, 点赞数：{like})")
                        print(f"       {msg}")

                        emotes = reply["content"].get("emote")
                        if emotes is not None:
                            for emote in emotes:
                                if emote[0] == "[":
                                    ok = emojis.get(emote)
                                    if ok is not None:
                                        emojis[emote] += 1
                                    else:
                                        emojis[emote] = 1

                if exit:
                    break
                pn += 1

                time.sleep(1)

        if br:
            break
        pa += 1
        time.sleep(1)

    print("表情信息：")
    print(emojis)


if __name__ == "__main__":
    bvid = input("输入BV号:")
    oid = bv2av(bvid=bvid)  # 评论区oid好像就是av号
    main()
