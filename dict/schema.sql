create database ldoce;
use ldoce;

create table word
(
    id      int(11)      NOT NULL AUTO_INCREMENT,
    word    varchar(255) NOT NULL,
    content longtext     NOT NULL,
    ref_word longtext NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_word_unique (word)
);

create table word_family
(
    id            int(11)      NOT NULL AUTO_INCREMENT,
    word_id       int(11)      NOT NULL,
    word_family   varchar(255) NULL,
    pronunciation varchar(40)  NULL,
    resource      json         NULL,
    PRIMARY KEY (id)
);

create table meaning
(
    id               int(11)  NOT NULL AUTO_INCREMENT,
    word_id          int(11)  NOT NULL,
    word_family_id   int(11)  NOT NULL,
    meaning          longtext NOT NULL,
    meaning_html     longtext NOT NULL,
    meaning_category varchar(50),
    position         int(11),
    meaning_type    varchar(20) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (word_family_id) REFERENCES word_family (id),
    FOREIGN KEY (word_id) REFERENCES word (id)
)