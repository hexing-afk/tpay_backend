(function(q,l){var f=function(a,c,b,d){"undefined"!=typeof console&&console&&console.warn&&console.warn("[ERROR]","{"+a+"}","<"+c+">",'"'+(b||"")+'"',d||"")},k=function(a,c,b,d){/[?&]upadm_debug(=|&|$)/.test(location.href)&&"undefined"!=typeof console&&console&&console.log&&console.log("[DEBUG]","{"+a+"}","<"+c+">",'"'+(b||"")+'"',d||"")},r,s=function(a){var c=[],b;for(b in a)c.push(b+"="+a[b]);return c.join("&")},g=function(a,c){for(var b in c)a[b]=c[b];return a},t=function(a){for(var c={locationId:a.locationId,
type:a.locationInfo.type?a.locationInfo.type:"default",tmpl:a.locationInfo.tmpl?a.locationInfo.tmpl:null},b=[],d=0,e=a.itemInfos.length;d<e;d++)b.push({itemId:a.itemInfos[d].itemId,name:a.itemInfos[d].name,title:a.itemInfos[d].title,href:a.itemInfos[d].href,width:a.itemInfos[d].width,height:a.itemInfos[d].height,content:a.itemInfos[d].content,type:a.itemInfos[d].type?a.itemInfos[d].type:c.type,tmpl:a.itemInfos[d].tmpl?a.itemInfos[d].tmpl:null,exposure:a.itemInfos[d].exposure});return{locationData:c,
itemDataList:b}},m=function(a,c){var b;if(a.call)b=a.call(null,c);else{b=a;for(var d in c)b=b.replace(RegExp("{"+d+"}","g"),c[d])}return b},u=function(a){if(1<a.itemDataList.length){for(var c=[],b=0,d=a.itemDataList.length;b<d;b++){var e=a.itemDataList[b].href,e=""!=e?a.itemDataList[b].tmpl||UPADM.tmpl_item[a.itemDataList[b].type]||UPADM.tmpl_item["default"]:a.itemDataList[b].tmpl||UPADM.tmpl_no_href_item[a.itemDataList[b].type]||UPADM.tmpl_item["default"];c.push('<li class="'+h.ITEM+'">'+m(e,a.itemDataList[b])+
"</li>")}a=a.locationData.tmpl||UPADM.tmpl_loc[a.locationData.type]||UPADM.tmpl_loc["default"];a='<div class="'+h.LOCATION+'">'+m(a,{content:'<ul class="'+h.ITEM_LIST+'">'+c.join("")+"</ul>"})+"</div>"}else e=a.itemDataList[0].href,e=""!=e?a.itemDataList[0].tmpl||UPADM.tmpl_item[a.itemDataList[0].type]||UPADM.tmpl_item["default"]:a.itemDataList[0].tmpl||UPADM.tmpl_no_href_item[a.itemDataList[0].type]||UPADM.tmpl_item["default"],c=m(e,a.itemDataList[0]),a=a.locationData.tmpl||UPADM.tmpl_loc[a.locationData.type]||
UPADM.tmpl_loc["default"],a='<div class="'+h.LOCATION+'">'+m(a,{content:c})+"</div>";return a},v=function(a,c){k("render","render","start",a);if(c)document.write(u(a));else{var b=l.getElementById(n+a.locationData.locationId);b?(b.innerHTML=u(a),k("render","render","end",a)):f("render","render","container not found",a)}for(var b=0,d=a.itemDataList.length;b<d;b++)a.itemDataList[b].exposure&&y(a.itemDataList[b].exposure)},o=function(a,c){var b=l.getElementById(n+a);if(b)switch(b.setAttribute(w.UPGG_STATE,
c),b.className=b.className+" "+h.CONTAINER,c){case j.START:b.style.background="transparent url("+z+") no-repeat center center";UPADM.containers[a]=b;break;case j.DONE:b.style.background="transparent"}else f("markRenderState","","container missing",a)},x=function(a,c){k("handleResponse","handleResponse","start",a);if(p("OK")==a.code)try{c.call(null,a.data)}catch(b){f("handleResponse","response exception","",a.data)}else f("handleResponse","response error",a.code,a)},y=function(a){k("exposure","","",
a);var c=document.createElement("img");c.src=a+(-1==a.indexOf("?")?"?":"&")+"random="+(new Date).getTime();var b=document.getElementsByTagName("head")[0],d=!1;c.onerror=c.onload=c.onreadystatechange=function(){if(!d&&(!this.readyState||"loaded"==this.readyState||"complete"==this.readyState))d=!0,b.removeChild(c)};b.appendChild(c)},z="https://static.95516.com/static/basis/images/loading.gif",w={UPGG_LOC_ID:"data-upgg-location-id",UPGG_GEN:"data-upgg-generated",UPGG_STATE:"data-upgg-state"},h={CONTAINER:"upgg_container",
LOCATION:"upgg_location",ITEM_LIST:"upgg_item-list",ITEM:"upgg_item"},n="upgg_location-",j={START:"1",DONE:"9"},p,A={OK:"00"};p=function(a){return A[a]};r={containers:{},data:{},renderAtOnce:function(a,c){c=c||{};o(a,j.START);var b;b=c||{};b=g(b,{callback:"UPADM.callback_renderAtOnce",locationIds:a});b="https://www.95516.com/ads/ads/g.do?"+s(b)+"&_t="+(new Date).getTime();k("request_renderAtOnce","callback","start",{id:a,url:b});document.write('<script type="text/javascript" src="'+b+'"><\/script>')},
autoRender:function(a){for(var a=a||{},c=l.getElementsByTagName("*"),b=[],d=0,e=c.length;d<e;d++)if(c[d].id&&0==c[d].id.indexOf(n)){var f=c[d].id.substr(n.length),h=c[d].getAttribute(w.UPGG_STATE);f&&h<j.START&&(b.push(f),o(f,j.START))}a=a||{};a=g(a,{callback:"UPADM.callback_autoRender",locationIds:b.join(",")});a="https://www.95516.com/ads/ads/g.do?"+s(a)+"&_t="+(new Date).getTime();k("request_autoRender","callback","start",{ids:b,url:a});b=l.createElement("script");b.src=a;l.getElementsByTagName("head")[0].appendChild(b)},
callback_renderAtOnce:function(a){x(a,function(a){var b=a[0];p("OK")==b.code&&b.locationId&&b.locationInfo&&b.itemInfos?(a=t(b),v(a,!0),o(b.locationId,j.DONE)):f("callback_renderAtOnce","response item",a[i].code,a[i])})},callback_autoRender:function(a){x(a,function(c){for(var c=a.data,b=0,d=c.length;b<d;b++)if(p("OK")==c[b].code&&c[b].locationId&&c[b].locationInfo&&c[b].itemInfos){var e=t(c[b]);v(e);o(c[b].locationId,j.DONE);setTimeout(function(){$("#pos_notify").show();var a=$("#upopNoticeInfo").height();
$("#pos_notify").css("top",a+44+"px");$("#orderList").css("margin-top",a+$("#pos_notify").height()+"px")},300)}else $("#pos_notify").hide(),e=$("#upopNoticeInfo").height(),$("#pos_notify").css("top",e+44+"px"),$("#orderList").css("margin-top",e+$("#pos_notify").height()+"px"),f("callback_autoRender","response item",c[b].code,c[b])})},tmpl_loc:{"default":"{content}"},tmpl_item:{1:'                    <a class="upgg_1_img-link" target="_blank" href="{href}">                        <img src="{content}" alt="{title}" width="{width}" height="{height}">                    </a>                    ',
2:'                    <div class="upgg_2_img-wrap">                        <a target="_blank" href="{href}">                            <img src="{content}" alt="{title}" width="{width}" height="{height}">                        </a>                    </div>                    <div class="upgg_2_title-wrap">                        <a target="_blank" href="{href}">{title}</a>                    </div>                    ',3:'<a target="_blank" href="{href}">{title}</a>',"default":"{title}"},tmpl_no_href_item:{1:'<img src="{content}" alt="{title}" width="{width}" height="{height}">                    ',
2:'<div class="upgg_2_img-wrap">                        <img src="{content}" alt="{title}" width="{width}" height="{height}">                    </div>                    <div class="upgg_2_title-wrap">                        {title}                    </div>                    ',3:"{title}","default":"{title}"},exposure:function(){},extend:function(a){g(this.tmpl_loc,a.tmpl_loc);g(this.tmpl_item,a.tmpl_item);g(this.tmpl_no_href_item,a.tmpl_no_href_item);g(this.containers,a.containers);g(this.data,
a.data);return this}};q.UPADM?q.UPADM.extend(r):q.UPADM=r})(window,document);