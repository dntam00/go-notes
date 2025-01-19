//全局配置默认展开折叠, 点击后可以展开折叠
ldoce_hide_example = false;
ldoce_hide_thesaurus = true;
ldoce_hide_copus_example = true;
ldoce_hide_collocations = true;
ldoce_hide_grammar = true;
is_ldoce_loaded = 'undefined' != typeof(is_ldoce_loaded)
console.log("~~~~~~~is_ldoce_loaded ",is_ldoce_loaded)
// 隐藏中文翻译
ldoce_hide_zh_translation = false
// 显示多词义导航(仅有多词义显示)
ldoce_show_nav = true;
// 显示ACTIV
ldoce_show_ACTIV = false;

// 在线发音 (默认离线) 启用后,可以删除 1.mdd
ldoce_online_pron = false;
// 在线图片 (默认离线) 启用后,可以删除 2.mdd
ldoce_online_image = false;

ldoce_root_class = "leon-ldoce"
ldoce_root_selector = "." + ldoce_root_class
ldoce_nav_class = "nav-ldoce"
ldoce_nav_selector = "." + ldoce_nav_class
debug = false
showUa = debug && true
showUa && alert(navigator.userAgent)


var LDOCE_PREFIX_PRON = "https://www.ldoceonline.com/media/english/";
var LDOCE_PREFIX_IMAGE = "https://www.ldoceonline.com/media/english/illustration/";

var LEDOCE_DICT_MAPPING = {
    "From Longman Dictionary of Contemporary English": "LDoCE ",
    "From Longman Business Dictionary": "LBD ",
}
// 挂载全局对象,进行方法隔离
ldoce = {
    root: function () {
        return document.querySelector(ldoce_root_selector);
    },
    addNavigation: function () {
        var root = ldoce.root()
        if (!root) {
            return;
        }
        // 已添加不添加
        if (root.querySelector(ldoce_nav_selector))
            return;

        let eles = root.querySelectorAll(".dictentry");
        // 只有一个entry不添加
        if (eles.length < 2)
            return;
        let container = document.createElement('div');
        container.setAttribute("class", ldoce_nav_class);
        for (let i = 0; i < eles.length; i++) {
            let sp = document.createElement('span');
            let pos = eles[i].querySelector(".POS");
            let dict = eles[i].querySelector(".dictionary_intro");
            pos = pos && pos.textContent;
            pos = !pos ? eles[i].querySelector(".HYPHENATION").textContent : pos;
            sp.textContent = (dict ? LEDOCE_DICT_MAPPING[dict.textContent] : "") + pos;
            if (i === 0) {
                sp.classList.add('active');
            }
            sp.addEventListener("click", function () {
                let index = ldoce.navIndex(this);
                console.log(index);
                for (let j = 0; j < container.children.length; j++) {
                    container.children[j].classList.remove('active');
                }
                this.classList.add('active');
                ldoce.showHideEntry(index)
            })
            container.appendChild(sp);
        }
        // 新增All
        let sp = document.createElement('span');
        sp.textContent = "All";
        container.appendChild(sp)
        sp.addEventListener("click", function () {
            for (let j = 0; j < container.children.length; j++) {
                container.children[j].classList.remove('active');
            }
            this.classList.add('active');
            ldoce.showHideEntry(-1)
        })
        ldoce.showHideEntry(0);
        root.insertBefore(container, root.childNodes[0]);
    },
    navIndex: function (ele) {
        var eles = document.querySelectorAll(ldoce_nav_selector + " span");
        for (var i = 0; i < eles.length; i++) {
            if (eles[i] === ele)
                return i;
        }
    },
    showHideEntry: function (index) {
        var root = ldoce.root();
        var eles = root.querySelectorAll(".dictentry");
        for (var i = 0; i < eles.length; i++) {
            if (index === i || index < 0)
                showEle(eles[i]);
            else
                hideEle(eles[i]);
        }
    },
    get_online_pron_url: function (src) {
        var parts = src.split('/');
        var dir = parts[parts.length - 2];
        var name = parts[parts.length - 1];
        return LDOCE_PREFIX_PRON + dir + "/" + name;
    },

    get_online_image_url: function (src) {
        var parts = src.split('/');
        var name = parts[parts.length - 1];
        return LDOCE_PREFIX_IMAGE + name;
    },
    nodeDisplay: function (node, value, fromThis, includeThis, excludeFunc, inlineList) {
        var nodes = node.parentElement.childNodes;
        var isFindThis = false
        for (var i = 0; i < nodes.length; i++) {
            let currentNode = nodes[i];
            if (excludeFunc && excludeFunc(currentNode)) {
                continue;
            }
            if (fromThis) {
                if (node === currentNode) {
                    isFindThis = true;
                }
                if (!isFindThis)
                    continue;
            }
            if (!includeThis && node === currentNode)
                continue;

            if (currentNode.style) {
                var display = "block"
                if (inlineList && inlineList.indexOf(currentNode.className) > -1) {
                    display = "inline"
                }
                currentNode.style.display = value === "block" ? display : value;
            }
        }
    }
}

!is_ldoce_loaded && document.addEventListener('DOMContentLoaded', function () {
    var eles = get_elements(".crossRef", function (ele) {
        return ele.textContent.indexOf("Verb table") > -1
    });
    for (var i = 0; i < eles.length; i++) {
        var url = eles[i].href;
        // 动词形态当前页面跳转,不重新加载新链接
        eles[i].href =
            url.indexOf("gdanchor=") > -1 ?
                "#" + url.substring(url.indexOf("=") + 1) : url.substring(url.indexOf("#"));

        eles[i].addEventListener("click", function (event) {
            document.querySelector(".verbTable").style.display = "block";
        })
    }
    elementsFn(".next_tenses", function (ele) {
        ele.style.display = "table-row"
    })


    var excludeClassName = ["SYN", "OPP", "BREQUIV", "RELATEDWD"]
    elementsClick(".DEF", function (ele) {
        if (!ele.nextElementSibling)
            return;

        let next = ele.nextElementSibling;
        while (next && excludeClassName.indexOf(next.className) > -1) {
            next = next.nextElementSibling
        }
        var state = next && next.style && next.style.display;

        if (state === "none") {
            ldoce.nodeDisplay(ele, "block", true, false, function (ele) {
                return excludeClassName.indexOf(ele.className) > -1
            });
            ele.classList.add('expanded');
        } else {
            ldoce.nodeDisplay(ele, "none", true, false, function (ele) {
                return excludeClassName.indexOf(ele.className) > -1
            });
            ele.classList.remove('expanded');
        }
    })

    ldoce_hide_example && elementsHide(".DEF~span");
    if (!ldoce_hide_example) {
        elementsFn(".DEF", function (ele) {
            ele.classList.add('expanded');
        });
    }
    ldoce_hide_thesaurus && elementsHide(".ThesBox>span:not(:first-child)");
    ldoce_hide_copus_example && elementsHide(".exaGroup>span:not(:first-child)");
    ldoce_hide_collocations && elementsHide(".ColloBox>span:not(:first-child)");
    ldoce_hide_grammar && elementsHide(".GramBox>span:not(:first-child)");
    elementsClick(".dictionary_intro,.F2NBox>span:first-child,.ThesBox>span:first-child,.ColloBox>span:first-child,.GramBox>span:first-child,.exaGroup>span:first-child", function (ele) {
        if (!ele.nextElementSibling)
            return;
        var state = ele.nextElementSibling.style.display;
        console.log("state tg", state, !state)
        if (state === "none") {
            ldoce.nodeDisplay(ele, "block", true);
            ele.classList.add('expanded');
        } else {
            ldoce.nodeDisplay(ele, "none", true);
            ele.classList.remove('expanded');
        }
    });

    ldoce_online_image && elementsFn("img", function (element) {
        console.log(element);
        element.src = ldoce.get_online_image_url(element.src);
        console.log(element);
    });
    ldoce_online_pron && elementsFn(".fa-volume-up", function (element) {
        console.log(element);
        element.href = ldoce.get_online_pron_url(element.href);
        element.onclick = function (e) {
            e.preventDefault();
            new Audio(element.href).play();
        }
    });

    ldoce_show_nav && ldoce.addNavigation();

    // 词头点击进行音节显示切换
    elementsClick(".HYPHENATION", function (ele) {
        if (ele.hasAttribute("syllable")) {
            ele.textContent = ele.textContent.indexOf("‧") > -1 ?
                ele.textContent.replaceAll("‧", "")
                : ele.getAttribute("syllable")
        } else if (ele.textContent.indexOf("‧") > -1) {
            ele.setAttribute("syllable", ele.textContent)
            ele.textContent = ele.textContent.replaceAll("‧", "");
        }
    })
    ldoce_hide_zh_translation && elementsHide(".cn");
    ldoce_show_ACTIV && elementsDisplay(".ldoceEntry .ACTIV,.ldoceEntry .ACTIV::before", "inline");
    debug && vConsole();
});

// 通用

function elementsFn(selector, fn) {
    var eles = document.querySelectorAll(selector);
    for (var i = 0; i < eles.length; i++) {
        fn(eles[i]);
    }
}

function elementsClick(selector, fn) {
    elementsFn(selector, function (ele) {
        ele.addEventListener("click", function (event) {
            fn(ele, event);
        })
    })
}

function elementsDisplay(selector, value) {
    elementsFn(selector, function (ele) {
        ele.style.display = value;
    })
}

function elementsHide(selector) {
    elementsDisplay(selector, "none")
}

function elementsShow(selector) {
    elementsDisplay(selector, "block")
}

function get_elements() {
    var selector = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : "input";
    var cond = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : function (el) {
        return el;
    };
    var eles = document.querySelectorAll(selector);
    var elements = [];
    for (var i = 0; i < eles.length; i++) {
        if (cond(eles[i])) {
            elements.push(eles[i]);
        }
    }
    return elements;
}

function hideEle(ele) {
    ele.style.display = 'none';
}

function showEle(ele) {
    ele.style.display = 'block';
}
function addScript(src, func) {
    let script = document.createElement('script');
    script.type = 'text/javascript';
    document.body.appendChild(script);
    if (src) {
        script.src = src;
    }
    script.onload = func;
}

function vConsole() {
    if ('undefined' == typeof(VConsole)) {
        addScript("https://cdn.bootcss.com/vConsole/3.15.0/vconsole.min.js", function () {
            var vConsole = new VConsole();
        })
    }
}