package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func main() {
	testTailRefEntry()
}

func testWordMeaning() {
	word := "bandwagon"
	data := `
		<link rel="stylesheet" href="ldoce.css"> <script type="text/javascript" src="ldoce.js"></script><div class="content leon-ldoce"> <div>  <div class="entry_content"><h1 class="pagetitle">bandwagon</h1><div class="dictionary"> <span class="dictentry"><span class="dictionary_intro span">From Longman Dictionary of Contemporary English</span><span class="dictlink"><span class="ldoceEntry Entry"><span class="Head"><span class="HWD">bandwagon</span><span class="HYPHENATION">band‧wag‧on</span><span class="PronCodes"><span class="neutral span"> /</span><span class="PRON">ˈbændˌwæɡən</span><span class="neutral span">/</span></span><span class="POS"> noun</span><span class="GRAM"><span class="neutral span"> [</span>countable<span class="neutral span">]</span></span> <a href="sound://breProns/bandwagon0205.mp3" class="brefile fas fa-volume-up hideOnAmp" title="Play British pronunciation of bandwagon">&nbsp;</a> <a href="sound://ameProns/bandwagon.mp3" class="amefile fas fa-volume-up hideOnAmp" title="Play American pronunciation of bandwagon">&nbsp;</a></span><span class="Sense" id="bandwagon__1"><span class="sensenum span">1</span> <span class="ACTIV">TAKE PART/BE INVOLVED</span><span class="DEF">an <a class="defRef" title="activity" href="entry://activity">activity</a> that a lot of people are doing</span><span class="EXAMPLE"> <a href="sound://exaProns/p008-000933297.mp3" class="exafile fas fa-volume-up hideOnAmp" title="Play Example">&nbsp;</a>The keep-fit bandwagon started rolling in the mid-80s.</span></span><span class="Sense" id="bandwagon__2"><span class="sensenum span">2</span> <span class="Crossref"><span class="neutral span"> →</span><a title="climb/jump/get on the bandwagon" class="crossRef" href="entry://climb-jump-get-on-the-bandwagon"> <span class="REFHWD">climb/jump/get on the bandwagon</span></a></span></span></span></span><span class="asset div"><span class="yellow_box"><span class="asset_intro">Examples from the Corpus</span></span></span><span class="assetlink"><span class="exaGroup cexa1 exaGroup"><span class="title">bandwagon</span><span class="cexa1g1 exa"><span class="neutral span">• </span>There is a <span class="NodeW">bandwagon</span> effect that is <a class="defRef" title="apparent" href="entry://apparent">apparent</a> once <a class="defRef" title="initiative" href="entry://initiative">initiatives</a> are taken.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>As the J-Boat <span class="NodeW">bandwagon</span> <a class="defRef" title="gather" href="entry://gather">gathered</a> <a class="defRef" title="momentum" href="entry://momentum">momentum</a>, other <a class="defRef" title="design" href="entry://design">designs</a> took <a class="defRef" title="shape" href="entry://shape">shape</a> on <a class="defRef" title="rod" href="entry://rod">Rod</a> Johnstone's <a class="defRef" title="board" href="entry://board">board</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>In every country, <a class="defRef" title="intellectual" href="entry://intellectual">intellectuals</a>, too, have <a class="defRef" title="jump" href="entry://jump">jumped</a> on the <a class="defRef" title="nationalist" href="entry://nationalist">nationalist</a> <span class="NodeW">bandwagon</span>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>So how do you <a class="defRef" title="hop" href="entry://hop">hop</a> on the <span class="NodeW">bandwagon</span>?</span><span class="cexa1g1 exa"><span class="neutral span">• </span>Just a <a class="defRef" title="preliminary" href="entry://preliminary">preliminary</a> <a class="defRef" title="communication" href="entry://communication">communication</a> first, without the <a class="defRef" title="experimental" href="entry://experimental">experimental</a> <a class="defRef" title="detail" href="entry://detail">details</a>, so that <a class="defRef" title="nobody" href="entry://nobody#nobody__3">nobody</a> can jump on the <span class="NodeW">bandwagon</span> right away.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>Companies such as <a class="defRef" title="oracle" href="entry://oracle">Oracle</a> are jumping on the <span class="NodeW">bandwagon</span>, too, with low-priced <a class="defRef" title="network" href="entry://network">network</a> computers.</span></span></span></span> </div> </div> </div>  </div>
	`
	reader := strings.NewReader(data)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("create document failed", err)
	}

	find := document.Find(".ldoceEntry.Entry")
	find.Each(func(i int, s *goquery.Selection) {
		idSelector := fmt.Sprintf("[id^=\"%s__\"]", word)
		s.Find(idSelector).Each(func(i int, s *goquery.Selection) {
			defHtml, _ := s.Html()
			fmt.Printf("html definition %d: %s\n", i, defHtml)

			//word1 := s.Text()
			s.Find(".DEF").Each(func(i int, a *goquery.Selection) {
				fmt.Printf("text definition: %s\n", a.Text())
			})

			s.Find(".crossRef").Each(func(i int, a *goquery.Selection) {
				fmt.Printf("crossRef: %s\n", a.Text())
				val, _ := a.Attr("href")
				fmt.Printf("attr: %s\n", strings.TrimPrefix(val, "entry://"))
			})
			//fmt.Printf("Word %d: %s\n", i, word1)
		})
	})
}

func testTailRefEntry() {
	data := `<link rel="stylesheet" href="ldoce.css"> <script type="text/javascript" src="ldoce.js"></script><div class="content leon-ldoce"> <div>  <div class="entry_content"><h1 class="pagetitle">magnesia</h1><div class="dictionary"> <span class="dictentry"><span class="dictionary_intro span">From Longman Dictionary of Contemporary English</span><span class="dictlink"><span class="ldoceEntry Entry"><span class="topics_container"><span class="related_topics">Related topics: </span><a class="topic" title="Drugs, medicines topic" href="entry://topic_Drugs, medicines">Drugs, medicines</a></span><span class="Head"><span class="HWD">magnesia</span><span class="HYPHENATION">mag‧ne‧sia</span><span class="PronCodes"><span class="neutral span"> /</span><span class="PRON">mæɡˈniːʃə, -ʒə</span><span class="neutral span">/</span></span><span class="POS"> noun</span><span class="GRAM"><span class="neutral span"> [</span>uncountable<span class="neutral span">]</span></span> <a href="sound://breProns/ld41magnesia.mp3" class="brefile fas fa-volume-up hideOnAmp" title="Play British pronunciation of magnesia">&nbsp;</a> <a href="sound://ameProns/magnesia.mp3" class="amefile fas fa-volume-up hideOnAmp" title="Play American pronunciation of magnesia">&nbsp;</a></span><span class="Sense" id="magnesia__1"><span class="FIELD">HCC</span><span class="FIELD">MD</span></span><span class="Tail"><span class="Crossref"><span class="neutral span"> →</span><a title="milk of magnesia" class="crossRef" href="entry://milk-of-magnesia"> <span class="REFHWD">milk of magnesia</span></a></span></span></span></span><span class="asset div"><span class="yellow_box"><span class="asset_intro">Examples from the Corpus</span></span></span><span class="assetlink"><span class="exaGroup cexa1 exaGroup"><span class="title">magnesia</span><span class="cexa1g1 exa"><span class="neutral span">• </span>They <a class="defRef" title="contain" href="entry://contain">contain</a> between 0.6 % and 0.8 % <span class="NodeW">magnesia</span> and contain <a class="defRef" title="low" href="entry://low">low</a> <a class="defRef" title="potassium" href="entry://potassium">potassium</a> <a class="defRef" title="oxide" href="entry://oxide">oxide</a> <a class="defRef" title="level" href="entry://level">levels</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>The <a class="defRef" title="clay" href="entry://clay">clays</a> <a class="defRef" title="consist" href="entry://consist">consist</a> of <a class="defRef" title="silica" href="entry://silica">silica</a> tetrahedra and the octahedra contain <span class="NodeW">magnesia</span> <a class="defRef" title="surround" href="entry://surround">surrounded</a> by <a class="defRef" title="oxygen" href="entry://oxygen">oxygen</a> <a class="defRef" title="atom" href="entry://atom">atoms</a> and hydroxyl groups.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>Some <a class="defRef" title="glaze" href="entry://glaze">glazes</a> contain <a class="defRef" title="elevate" href="entry://elevate">elevated</a> <span class="NodeW">magnesia</span> which increases their durability; this was certainly <a class="defRef" title="necessary" href="entry://necessary">necessary</a> given the low <a class="defRef" title="calcium" href="entry://calcium">calcium</a> oxide levels.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>Other <a class="defRef" title="analysis" href="entry://analysis">analyses</a> of soda-rich <a class="defRef" title="plant" href="entry://plant">plant</a> <a class="defRef" title="ash" href="entry://ash">ashes</a> contained considerably higher <span class="NodeW">magnesia</span> levels <a class="defRef" title="accompany" href="entry://accompany">accompanying</a> the <a class="defRef" title="soda" href="entry://soda">soda</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span><a class="defRef" title="milk" href="entry://milk">Milk</a> of <span class="NodeW">magnesia</span>, an osmotic <a class="defRef" title="laxative" href="entry://laxative">laxative</a>, was used <a class="defRef" title="accord" href="entry://accord">according</a> to <a class="defRef" title="age" href="entry://age">age</a>, body <a class="defRef" title="weight" href="entry://weight">weight</a>, and severity of the <a class="defRef" title="constipation" href="entry://constipation">constipation</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>The compositional <a class="defRef" title="difference" href="entry://difference">differences</a> are principally <a class="defRef" title="due" href="entry://due">due</a> to <a class="defRef" title="impurity" href="entry://impurity">impurity</a> levels of the oxides of <span class="NodeW">magnesia</span>, and lead in the <a class="defRef" title="glass" href="entry://glass">glasses</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span><a class="defRef" title="indigestion" href="entry://indigestion">Indigestion</a> could be <a class="defRef" title="quell" href="entry://quell">quelled</a> with a <a class="defRef" title="simple" href="entry://simple">simple</a> <span class="NodeW">magnesia</span> <a class="defRef" title="tablet" href="entry://tablet">tablet</a>, thousands of which she <a class="defRef" title="chew" href="entry://chew">chewed</a> in her <a class="defRef" title="lifetime" href="entry://lifetime">lifetime</a>.</span><span class="cexa1g1 exa"><span class="neutral span">• </span>The different soda, <span class="NodeW">magnesia</span> and <a class="defRef" title="phosphorus" href="entry://phosphorus">phosphorus</a> pentoxide levels can be <a class="defRef" title="relate" href="entry://relate">related</a> to the use of a different soda <a class="defRef" title="source" href="entry://source">source</a>.</span></span></span></span><span class="etym"><span class="asset_intro">Origin</span> <span class="Head"><span class="HWD">magnesia</span></span> <span class="Sense" id="magnesia__3"><span class="CENTURY"><span class="neutral span">(</span>1300-1400<span class="neutral span">)</span></span> <span class="LANG">Modern Latin</span> <span class="ORIGIN">magnes carneus</span> <span class="TRAN"><span class="neutral span">“</span>flesh magnet<span class="neutral span">”</span></span>, used of a white powder that stuck to the lips; <span class="neutral span"> →&nbsp;</span><a title="MAGNET" class="crossRef" href="entry://magnet#magnet__3"><span class="REFHWD">MAGNET</span></a></span></span> </div> </div> </div>  </div>`
	reader := strings.NewReader(data)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("create document failed", err)
	}

	refEntrySelection := document.Find(".Tail .Crossref .crossRef")
	refEntrySelection.Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		val, _ := s.Attr("href")
		fmt.Println(strings.TrimPrefix(val, "entry://"))
	})
}

func testRefEntry() {
	data := `
<link rel="stylesheet" href="ldoce.css">
<script type="text/javascript" src="ldoce.js"></script>
<div class="ref_entry">
    <div class="ref_head">advertere</div>
    <div class="ref_body">
        <div class="ref_item"><span class="ref_bullet">→</span><a href="entry://adverse#adverse__5">adverse</a></div>
        <div class="ref_item"><span class="ref_bullet">→</span><a href="entry://advert">advert</a></div>
        <div class="ref_item"><span class="ref_bullet">→</span><a href="entry://inadvertently">inadvertently</a></div>
    </div>
</div>`
	reader := strings.NewReader(data)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("create document failed", err)
	}

	// 2. reference to pure form of word
	refEntrySelection := document.Find(".ref_entry .ref_item a")
	refEntrySelection.Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
	if len(refEntrySelection.Nodes) > 0 {
		refWord := refEntrySelection.Text()
		refWord = strings.TrimSpace(refWord)
		fmt.Println(refWord)
	}
}
