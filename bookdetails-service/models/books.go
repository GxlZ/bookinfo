package models

import (
	"time"
)

type Books struct {
	ID          uint       `gorm:"primary_key"`
	Name        string     `gorm:"index:name"`   //书名
	Author      string                           //作者
	Intro       string     `gorm:"type:text"`    //简介
	PublishDate string     `json:"publish_date"` //出版时间
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

var books = []Books{
	{
		ID:          1,
		Name:        "西游记",
		Author:      "吴承恩",
		PublishDate: "明代",
		Intro:       "《西游记》是中国古代第一部浪漫主义章回体长篇神魔小说。现存明刊百回本《西游记》均无作者署名。清代学者吴玉搢等首先提出《西游记》作者是明代吴承恩 [1]  。这部小说以“唐僧取经”这一历史事件为蓝本，通过作者的艺术加工，深刻地描绘了当时的社会现实。全书主要描写了孙悟空出世及大闹天宫后，遇见了唐僧、猪八戒和沙僧三人，西行取经，一路降妖伏魔，经历了九九八十一难，终于到达西天见到如来佛祖，最终五圣成真的故事。",
	},
	{
		ID:          2,
		Name:        "水浒传",
		Author:      "施耐庵著,罗贯中整理",
		PublishDate: "明代",
		Intro:       "《水浒传》是一部长篇英雄传奇，是中国古代长篇小说的代表作之一，是以宋江起义故事为线索创作出来的",
	},
	{
		ID:          3,
		Name:        "三国演义",
		Author:      "罗贯中",
		PublishDate: "元末明初",
		Intro:       "《三国演义》是中国古典四大名著之一，是中国第一部长篇章回体历史演义小说，全名为《三国志通俗演义》（又称《三国志演义》），作者是元末明初的著名小说家罗贯中。《三国志通俗演义》成书后有嘉靖壬午本等多个版本传于世，到了明末清初，毛宗岗对《三国演义》整顿回目、修正文辞、改换诗文",
	},
}
