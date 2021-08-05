(function(){function l(a,b){this.wrapper="string"==typeof a?document.querySelector(a):a;this.scroller=this.wrapper.children[0];this.scrollerStyle=this.scroller.style;this.options={disablePointer:!0,disableTouch:!d.hasTouch,disableMouse:d.hasTouch,startX:0,startY:0,scrollY:!0,directionLockThreshold:5,momentum:!0,bounce:!0,bounceTime:600,bounceEasing:"",preventDefault:!0,preventDefaultException:{tagName:/^(INPUT|TEXTAREA|BUTTON|SELECT)$/},HWCompositing:!0,useTransition:!0,useTransform:!0,bindToWrapper:"undefined"===
typeof window.onmousedown};for(var c in b)this.options[c]=b[c];this.translateZ=this.options.HWCompositing&&d.hasPerspective?" translateZ(0)":"";this.options.useTransition=d.hasTransition&&this.options.useTransition;this.options.useTransform=d.hasTransform&&this.options.useTransform;this.options.eventPassthrough=!0===this.options.eventPassthrough?"vertical":this.options.eventPassthrough;this.options.preventDefault=!this.options.eventPassthrough&&this.options.preventDefault;this.options.scrollY="vertical"==
this.options.eventPassthrough?!1:this.options.scrollY;this.options.scrollX="horizontal"==this.options.eventPassthrough?!1:this.options.scrollX;this.options.freeScroll=this.options.freeScroll&&!this.options.eventPassthrough;this.options.directionLockThreshold=this.options.eventPassthrough?0:this.options.directionLockThreshold;this.options.bounceEasing="string"==typeof this.options.bounceEasing?d.ease[this.options.bounceEasing]||d.ease.circular:this.options.bounceEasing;this.options.resizePolling=void 0===
this.options.resizePolling?60:this.options.resizePolling;!0===this.options.tap&&(this.options.tap="tap");!this.options.useTransition&&!this.options.useTransform&&!/relative|absolute/i.test(this.scrollerStyle.position)&&(this.scrollerStyle.position="relative");"scale"==this.options.shrinkScrollbars&&(this.options.useTransition=!1);this.options.invertWheelDirection=this.options.invertWheelDirection?-1:1;3==this.options.probeType&&(this.options.useTransition=!1);this.directionY=this.directionX=this.y=
this.x=0;this._events={};this._init();this.refresh();this.scrollTo(this.options.startX,this.options.startY);this.enable()}function r(a,b){if(!(this instanceof r))return new r(a,b);this.html=a;this.opts=b;var c=document.createElement("div");c.className="olay";var f=document.createElement("div");f.className="layer";this.el=c;this.layer_el=f;this.init()}function s(a,b,c){if(!n.isArray(b)||0===b.length)throw new TypeError("the data must be a non-empty array!");if(-1==[1,2,3,4,5,6].indexOf(a))throw new RangeError("the level parameter must be one of 1,2,3,4,5,6!");
this.data=b;this.level=a||1;this.options=c;this.typeBox="one-level-box";a="one two three four five six".split(" ");6>=this.level&&1<=this.level&&(this.typeBox=a[parseInt(this.level)-1]+"-level-box");this.title=c.title||"";this.options.itemHeight=c.itemHeight||35;this.options.itemShowCount=-1!==[3,5,7,9].indexOf(c.itemShowCount)?c.itemShowCount:7;this.options.coverArea1Top=Math.floor(this.options.itemShowCount/2);this.options.coverArea2Top=Math.ceil(this.options.itemShowCount/2);this.options.headerHeight=
c.headerHeight||44;this.options.relation=n.isArray(this.options.relation)?this.options.relation:[];this.options.oneTwoRelation=this.options.relation[0];this.options.twoThreeRelation=this.options.relation[1];this.options.threeFourRelation=this.options.relation[2];this.options.fourFiveRelation=this.options.relation[3];this.options.fiveSixRelation=this.options.relation[4];"px"!==this.options.cssUnit&&"rem"!==this.options.cssUnit&&(this.options.cssUnit="px");this.selectOneObj={id:this.options.oneLevelId};
this.selectTwoObj={id:this.options.twoLevelId};this.selectThreeObj={id:this.options.threeLevelId};this.selectFourObj={id:this.options.fourLevelId};this.selectFiveObj={id:this.options.fiveLevelId};this.selectSixObj={id:this.options.sixLevelId};this.setBase();this.init()}var w=window.requestAnimationFrame||window.webkitRequestAnimationFrame||window.mozRequestAnimationFrame||window.oRequestAnimationFrame||window.msRequestAnimationFrame||function(a){window.setTimeout(a,1E3/60)},u;a:{for(var p=document.createElement("div").style,
h=["a","webkitA","MozA","OA","msA"],o=["animationend","webkitAnimationEnd","animationend","oAnimationEnd","MSAnimationEnd"],t,j=0,v=h.length;j<v;j++)if(t=h[j]+"nimation",t in p){u=o[j];break a}u="animationend"}var d,p=function(a){return!1===q?!1:""===q?a:q+a.charAt(0).toUpperCase()+a.substr(1)},g={},h=document.createElement("div").style,q;a:{o=["t","webkitT","MozT","msT","OT"];j=0;for(v=o.length;j<v;j++)if(t=o[j]+"ransform",t in h){q=o[j].substr(0,o[j].length-1);break a}q=!1}g.getTime=Date.now||function(){return(new Date).getTime()};
g.extend=function(a,b){for(var c in b)a[c]=b[c]};g.addEvent=function(a,b,c,f){a.addEventListener(b,c,!!f)};g.removeEvent=function(a,b,c,f){a.removeEventListener(b,c,!!f)};g.prefixPointerEvent=function(a){return window.MSPointerEvent?"MSPointer"+a.charAt(7).toUpperCase()+a.substr(8):a};g.momentum=function(a,b,c,f,i,e){var b=a-b,c=Math.abs(b)/c,d,e=void 0===e?6E-4:e;d=a+c*c/(2*e)*(0>b?-1:1);e=c/e;d<f?(d=i?f-i/2.5*(c/8):f,b=Math.abs(d-a),e=b/c):0<d&&(d=i?i/2.5*(c/8):0,b=Math.abs(a)+d,e=b/c);return{destination:Math.round(d),
duration:e}};o=p("transform");g.extend(g,{hasTransform:!1!==o,hasPerspective:p("perspective")in h,hasTouch:"ontouchstart"in window,hasPointer:!(!window.PointerEvent&&!window.MSPointerEvent),hasTransition:p("transition")in h});h=window.navigator.appVersion;h=/Android/.test(h)&&!/Chrome\/\d/.test(h)?(h=h.match(/Safari\/(\d+.\d)/))&&"object"===typeof h&&2<=h.length?535.19>parseFloat(h[1]):!0:!1;g.isBadAndroid=h;g.extend(g.style={},{transform:o,transitionTimingFunction:p("transitionTimingFunction"),transitionDuration:p("transitionDuration"),
transitionDelay:p("transitionDelay"),transformOrigin:p("transformOrigin")});g.hasClass=function(a,b){return RegExp("(^|\\s)"+b+"(\\s|$)").test(a.className)};g.addClass=function(a,b){if(!g.hasClass(a,b)){var c=a.className.split(" ");c.push(b);a.className=c.join(" ")}};g.removeClass=function(a,b){if(g.hasClass(a,b))a.className=a.className.replace(RegExp("(^|\\s)"+b+"(\\s|$)","g")," ")};g.offset=function(a){for(var b=-a.offsetLeft,c=-a.offsetTop;a=a.offsetParent;){b=b-a.offsetLeft;c=c-a.offsetTop}return{left:b,
top:c}};g.preventDefaultException=function(a,b){for(var c in b)if(b[c].test(a[c]))return true;return false};g.extend(g.eventType={},{touchstart:1,touchmove:1,touchend:1,mousedown:2,mousemove:2,mouseup:2,pointerdown:3,pointermove:3,pointerup:3,MSPointerDown:3,MSPointerMove:3,MSPointerUp:3});g.extend(g.ease={},{quadratic:{style:"cubic-bezier(0.25, 0.46, 0.45, 0.94)",fn:function(a){return a*(2-a)}},circular:{style:"cubic-bezier(0.1, 0.57, 0.1, 1)",fn:function(a){return Math.sqrt(1- --a*a)}},back:{style:"cubic-bezier(0.175, 0.885, 0.32, 1.275)",
fn:function(a){return(a=a-1)*a*(5*a+4)+1}},bounce:{style:"",fn:function(a){return(a=a/1)<1/2.75?7.5625*a*a:a<2/2.75?7.5625*(a=a-1.5/2.75)*a+0.75:a<2.5/2.75?7.5625*(a=a-2.25/2.75)*a+0.9375:7.5625*(a=a-2.625/2.75)*a+0.984375}},elastic:{style:"",fn:function(a){return a===0?0:a==1?1:0.4*Math.pow(2,-10*a)*Math.sin((a-0.055)*2*Math.PI/0.22)+1}}});g.tap=function(a,b){var c=document.createEvent("Event");c.initEvent(b,true,true);c.pageX=a.pageX;c.pageY=a.pageY;a.target.dispatchEvent(c)};g.click=function(a){var b=
a.target,c;if(!/(SELECT|INPUT|TEXTAREA)/i.test(b.tagName)){c=document.createEvent(window.MouseEvent?"MouseEvents":"Event");c.initEvent("click",true,true);c.view=a.view||window;c.detail=1;c.screenX=b.screenX||0;c.screenY=b.screenY||0;c.clientX=b.clientX||0;c.clientY=b.clientY||0;c.ctrlKey=!!a.ctrlKey;c.altKey=!!a.altKey;c.shiftKey=!!a.shiftKey;c.metaKey=!!a.metaKey;c.button=0;c.relatedTarget=null;c._constructed=true;b.dispatchEvent(c)}};d=g;l.prototype={version:"1.0.0",_init:function(){this._initEvents()},
destroy:function(){this._initEvents(true);clearTimeout(this.resizeTimeout);this.resizeTimeout=null;this._execEvent("destroy")},_transitionEnd:function(a){if(a.target==this.scroller&&this.isInTransition){this._transitionTime();if(!this.resetPosition(this.options.bounceTime)){this.isInTransition=false;this._execEvent("scrollEnd")}}},_start:function(a){if(!(d.eventType[a.type]!=1&&(a.which?a.button:a.button<2?0:a.button==4?1:2)!==0)&&this.enabled&&!(this.initiated&&d.eventType[a.type]!==this.initiated)){this.options.preventDefault&&
(!d.isBadAndroid&&!d.preventDefaultException(a.target,this.options.preventDefaultException))&&a.preventDefault();var b=a.touches?a.touches[0]:a;this.initiated=d.eventType[a.type];this.moved=false;this.directionLocked=this.directionY=this.directionX=this.distY=this.distX=0;this.startTime=d.getTime();if(this.options.useTransition&&this.isInTransition){this._transitionTime();this.isInTransition=false;a=this.getComputedPosition();this._translate(Math.round(a.x),Math.round(a.y));this._execEvent("scrollEnd")}else if(!this.options.useTransition&&
this.isAnimating){this.isAnimating=false;this._execEvent("scrollEnd")}this.startX=this.x;this.startY=this.y;this.absStartX=this.x;this.absStartY=this.y;this.pointX=b.pageX;this.pointY=b.pageY;this._execEvent("beforeScrollStart")}},_move:function(a){if(this.enabled&&d.eventType[a.type]===this.initiated){var b=a.touches?a.touches[0]:a,c=b.pageX-this.pointX,f=b.pageY-this.pointY,i=d.getTime(),e;this.pointX=b.pageX;this.pointY=b.pageY;this.distX=this.distX+c;this.distY=this.distY+f;b=Math.abs(this.distX);
e=Math.abs(this.distY);if(!(i-this.endTime>300&&b<10&&e<10)){if(!this.directionLocked&&!this.options.freeScroll)this.directionLocked=b>e+this.options.directionLockThreshold?"h":e>=b+this.options.directionLockThreshold?"v":"n";if(this.directionLocked=="h"){if(this.options.eventPassthrough=="vertical")a.preventDefault();else if(this.options.eventPassthrough=="horizontal"){this.initiated=false;return}f=0}else if(this.directionLocked=="v"){if(this.options.eventPassthrough=="horizontal")a.preventDefault();
else if(this.options.eventPassthrough=="vertical"){this.initiated=false;return}c=0}c=this.hasHorizontalScroll?c:0;f=this.hasVerticalScroll?f:0;a=this.x+c;b=this.y+f;if(a>0||a<this.maxScrollX)a=this.options.bounce?this.x+c/3:a>0?0:this.maxScrollX;if(b>0||b<this.maxScrollY)b=this.options.bounce?this.y+f/3:b>0?0:this.maxScrollY;this.directionX=c>0?-1:c<0?1:0;this.directionY=f>0?-1:f<0?1:0;this.moved||this._execEvent("scrollStart");this.moved=true;this._translate(a,b);if(i-this.startTime>300){this.startTime=
i;this.startX=this.x;this.startY=this.y;this.options.probeType==1&&this._execEvent("scroll")}this.options.probeType>1&&this._execEvent("scroll")}}},_end:function(a){if(this.enabled&&d.eventType[a.type]===this.initiated){this.options.preventDefault&&!d.preventDefaultException(a.target,this.options.preventDefaultException)&&a.preventDefault();var b,c;c=d.getTime()-this.startTime;var f=Math.round(this.x),i=Math.round(this.y),e=Math.abs(f-this.startX),g=Math.abs(i-this.startY);b=0;var k="";this.initiated=
this.isInTransition=0;this.endTime=d.getTime();if(!this.resetPosition(this.options.bounceTime)){this.scrollTo(f,i);if(this.moved)if(this._events.flick&&c<200&&e<100&&g<100)this._execEvent("flick");else{if(this.options.momentum&&c<300){b=this.hasHorizontalScroll?d.momentum(this.x,this.startX,c,this.maxScrollX,this.options.bounce?this.wrapperWidth:0,this.options.deceleration):{destination:f,duration:0};c=this.hasVerticalScroll?d.momentum(this.y,this.startY,c,this.maxScrollY,this.options.bounce?this.wrapperHeight:
0,this.options.deceleration):{destination:i,duration:0};f=b.destination;i=c.destination;b=Math.max(b.duration,c.duration);this.isInTransition=1}if(this.options.snap){this.currentPage=k=this._nearestSnap(f,i);b=this.options.snapSpeed||Math.max(Math.max(Math.min(Math.abs(f-k.x),1E3),Math.min(Math.abs(i-k.y),1E3)),300);f=k.x;i=k.y;this.directionY=this.directionX=0;k=this.options.bounceEasing}if(f!=this.x||i!=this.y){if(f>0||f<this.maxScrollX||i>0||i<this.maxScrollY)k=d.ease.quadratic;this.scrollTo(f,
i,b,k)}else this._execEvent("scrollEnd")}else{this.options.tap&&d.tap(a,this.options.tap);this.options.click&&d.click(a);this._execEvent("scrollCancel")}}}},_resize:function(){var a=this;clearTimeout(this.resizeTimeout);this.resizeTimeout=setTimeout(function(){a.refresh()},this.options.resizePolling)},resetPosition:function(a){var b=this.x,c=this.y;if(!this.hasHorizontalScroll||this.x>0)b=0;else if(this.x<this.maxScrollX)b=this.maxScrollX;if(!this.hasVerticalScroll||this.y>0)c=0;else if(this.y<this.maxScrollY)c=
this.maxScrollY;if(b==this.x&&c==this.y)return false;this.scrollTo(b,c,a||0,this.options.bounceEasing);return true},disable:function(){this.enabled=false},enable:function(){this.enabled=true},refresh:function(){this.wrapperWidth=this.wrapper.clientWidth;this.wrapperHeight=this.wrapper.clientHeight;this.scrollerWidth=this.scroller.offsetWidth;this.scrollerHeight=this.scroller.offsetHeight;this.maxScrollX=this.wrapperWidth-this.scrollerWidth;this.maxScrollY=this.wrapperHeight-this.scrollerHeight;this.hasHorizontalScroll=
this.options.scrollX&&this.maxScrollX<0;this.hasVerticalScroll=this.options.scrollY&&this.maxScrollY<0;if(!this.hasHorizontalScroll){this.maxScrollX=0;this.scrollerWidth=this.wrapperWidth}if(!this.hasVerticalScroll){this.maxScrollY=0;this.scrollerHeight=this.wrapperHeight}this.directionY=this.directionX=this.endTime=0;this.wrapperOffset=d.offset(this.wrapper);this._execEvent("refresh");this.resetPosition()},on:function(a,b){this._events[a]||(this._events[a]=[]);this._events[a].push(b)},off:function(a,
b){if(this._events[a]){var c=this._events[a].indexOf(b);c>-1&&this._events[a].splice(c,1)}},_execEvent:function(a){if(this._events[a]){var b=0,c=this._events[a].length;if(c)for(;b<c;b++)this._events[a][b].apply(this,[].slice.call(arguments,1))}},scrollTo:function(a,b,c,f){f=f||d.ease.circular;this.isInTransition=this.options.useTransition&&c>0;var i=this.options.useTransition&&f.style;if(!c||i){if(i){this._transitionTimingFunction(f.style);this._transitionTime(c)}this._translate(a,b)}else this._animate(a,
b,c,f.fn)},scrollToElement:function(a,b,c,f,i){if(a=a.nodeType?a:this.scroller.querySelector(a)){var e=d.offset(a);e.left=e.left-this.wrapperOffset.left;e.top=e.top-this.wrapperOffset.top;c===true&&(c=Math.round(a.offsetWidth/2-this.wrapper.offsetWidth/2));f===true&&(f=Math.round(a.offsetHeight/2-this.wrapper.offsetHeight/2));e.left=e.left-(c||0);e.top=e.top-(f||0);e.left=e.left>0?0:e.left<this.maxScrollX?this.maxScrollX:e.left;e.top=e.top>0?0:e.top<this.maxScrollY?this.maxScrollY:e.top;b=b===void 0||
b===null||b==="auto"?Math.max(Math.abs(this.x-e.left),Math.abs(this.y-e.top)):b;this.scrollTo(e.left,e.top,b,i)}},_transitionTime:function(a){if(this.options.useTransition){var a=a||0,b=d.style.transitionDuration;if(b){this.scrollerStyle[b]=a+"ms";if(!a&&d.isBadAndroid){this.scrollerStyle[b]="0.0001ms";var c=this;w(function(){c.scrollerStyle[b]==="0.0001ms"&&(c.scrollerStyle[b]="0s")})}}}},_transitionTimingFunction:function(a){this.scrollerStyle[d.style.transitionTimingFunction]=a},_translate:function(a,
b){if(this.options.useTransform)this.scrollerStyle[d.style.transform]="translate("+a+"px,"+b+"px)"+this.translateZ;else{a=Math.round(a);b=Math.round(b);this.scrollerStyle.left=a+"px";this.scrollerStyle.top=b+"px"}this.x=a;this.y=b},_initEvents:function(a){var a=a?d.removeEvent:d.addEvent,b=this.options.bindToWrapper?this.wrapper:window;a(window,"orientationchange",this);a(window,"resize",this);this.options.click&&a(this.wrapper,"click",this,true);if(!this.options.disableMouse){a(this.wrapper,"mousedown",
this);a(b,"mousemove",this);a(b,"mousecancel",this);a(b,"mouseup",this)}if(d.hasPointer&&!this.options.disablePointer){a(this.wrapper,d.prefixPointerEvent("pointerdown"),this);a(b,d.prefixPointerEvent("pointermove"),this);a(b,d.prefixPointerEvent("pointercancel"),this);a(b,d.prefixPointerEvent("pointerup"),this)}if(d.hasTouch&&!this.options.disableTouch){a(this.wrapper,"touchstart",this);a(b,"touchmove",this);a(b,"touchcancel",this);a(b,"touchend",this)}a(this.scroller,"transitionend",this);a(this.scroller,
"webkitTransitionEnd",this);a(this.scroller,"oTransitionEnd",this);a(this.scroller,"MSTransitionEnd",this)},getComputedPosition:function(){var a=window.getComputedStyle(this.scroller,null),b;if(this.options.useTransform){a=a[d.style.transform].split(")")[0].split(", ");b=+(a[12]||a[4]);a=+(a[13]||a[5])}else{b=+a.left.replace(/[^-\d.]/g,"");a=+a.top.replace(/[^-\d.]/g,"")}return{x:b,y:a}},_animate:function(a,b,c,f){function i(){var m=d.getTime(),j;if(m>=n){e.isAnimating=false;e._translate(a,b);e.resetPosition(e.options.bounceTime)||
e._execEvent("scrollEnd")}else{m=(m-h)/c;j=f(m);m=(a-g)*j+g;j=(b-k)*j+k;e._translate(m,j);e.isAnimating&&w(i);e.options.probeType==3&&e._execEvent("scroll")}}var e=this,g=this.x,k=this.y,h=d.getTime(),n=h+c;this.isAnimating=true;i()},handleEvent:function(a){switch(a.type){case "touchstart":case "pointerdown":case "MSPointerDown":case "mousedown":this._start(a);break;case "touchmove":case "pointermove":case "MSPointerMove":case "mousemove":this._move(a);break;case "touchend":case "pointerup":case "MSPointerUp":case "mouseup":case "touchcancel":case "pointercancel":case "MSPointerCancel":case "mousecancel":this._end(a);
break;case "orientationchange":case "resize":this._resize();break;case "transitionend":case "webkitTransitionEnd":case "oTransitionEnd":case "MSTransitionEnd":this._transitionEnd(a);break;case "click":if(this.enabled&&!a._constructed){a.preventDefault();a.stopPropagation()}}}};l.utils=d;var n={isArray:function(a){return Object.prototype.toString.call(a)==="[object Array]"},isFunction:function(a){return typeof a==="function"},attrToData:function(a,b){var c={},f;for(f in a.dataset)c[f]=a.dataset[f];
c.dom=a;c.atindex=b;return c},attrToHtml:function(a){var b="",c;for(c in a)b=b+("data-"+c+'="'+a[c]+'"');return b}};r.prototype={init:function(){this.layer_el.innerHTML=this.html;this.opts.container&&document.querySelector(this.opts.container)?document.querySelector(this.opts.container).appendChild(this.el):document.body.appendChild(this.el);this.el.appendChild(this.layer_el);this.el.style.height=Math.max(document.documentElement.getBoundingClientRect().height,window.innerHeight);if(this.opts.className)this.el.className=
this.el.className+(" "+this.opts.className);this.bindEvent()},bindEvent:function(){var a=this.el.querySelectorAll(".sure"),b=this.el.querySelectorAll(".close"),c=this;this.el.addEventListener("click",function(){c.close();c.opts.maskCallback&&c.opts.maskCallback()});this.layer_el.addEventListener("click",function(a){a.stopPropagation()});Array.prototype.slice.call(a).forEach(function(a){a.addEventListener("click",function(){c.close()})});Array.prototype.slice.call(b).forEach(function(a){a.addEventListener("click",
function(){c.close();c.opts.fallback&&c.opts.fallback()})})},close:function(){var a=this;if(a.el)if(a.opts.showAnimate){a.el.className=a.el.className+" fadeOutDown";a.el.addEventListener(u,function(){a.removeDom()})}else a.removeDom()},removeDom:function(){this.el.parentNode.removeChild(this.el);this.el=null;document.documentElement.classList.contains("ios-select-body-class")&&document.documentElement.classList.remove("ios-select-body-class")}};s.prototype={init:function(){this.initLayer();this.setLevelData(1,
this.options.oneLevelId,this.options.twoLevelId,this.options.threeLevelId,this.options.fourLevelId,this.options.fiveLevelId,this.options.sixLevelId)},initLayer:function(){var a=this,b=this.options.headerHeight+this.options.cssUnit,b=['<header style="height: '+b+"; line-height: "+b+'" class="iosselect-header">','<a style="height: '+b+"; line-height: "+b+'" href="javascript:void(0)" class="close">'+(this.options.closeText||"\u53d6\u6d88")+"</a>",'<a style="height: '+b+"; line-height: "+b+'" href="javascript:void(0)" class="sure">'+
(this.options.sureText||"\u786e\u5b9a")+"</a>",'<h2 id="iosSelectTitle"></h2>\r\n</header>\r\n<section class="iosselect-box">\r\n<div class="one-level-contain" id="oneLevelContain">\r\n<ul class="select-one-level">\r\n</ul>\r\n</div>\r\n<div class="two-level-contain" id="twoLevelContain">\r\n<ul class="select-two-level">\r\n</ul>\r\n</div>\r\n<div class="three-level-contain" id="threeLevelContain">\r\n<ul class="select-three-level">\r\n</ul>\r\n</div>\r\n<div class="four-level-contain" id="fourLevelContain">\r\n<ul class="select-four-level">\r\n</ul>\r\n</div>\r\n<div class="five-level-contain" id="fiveLevelContain">\r\n<ul class="select-five-level">\r\n</ul>\r\n</div>\r\n<div class="six-level-contain" id="sixLevelContain">\r\n<ul class="select-six-level">\r\n</ul>\r\n</div>\r\n</section>\r\n<hr class="cover-area1"/>\r\n<hr class="cover-area2"/>\r\n<div class="ios-select-loading-box" id="iosSelectLoadingBox">\r\n<div class="ios-select-loading"></div>\r\n</div>'].join("\r\n");
this.iosSelectLayer=new r(b,{className:"ios-select-widget-box "+this.typeBox+(this.options.addClassName?" "+this.options.addClassName:"")+(this.options.showAnimate?" fadeInUp":""),container:this.options.container||"",showAnimate:this.options.showAnimate,fallback:this.options.fallback,maskCallback:this.options.maskCallback});this.iosSelectTitleDom=this.iosSelectLayer.el.querySelector("#iosSelectTitle");this.iosSelectLoadingBoxDom=this.iosSelectLayer.el.querySelector("#iosSelectLoadingBox");this.iosSelectTitleDom.innerHTML=
this.title;if(this.options.headerHeight&&this.options.itemHeight){this.coverArea1Dom=this.iosSelectLayer.el.querySelector(".cover-area1");this.coverArea1Dom.style.top=this.options.headerHeight+this.options.itemHeight*this.options.coverArea1Top+this.options.cssUnit;this.coverArea2Dom=this.iosSelectLayer.el.querySelector(".cover-area2");this.coverArea2Dom.style.top=this.options.headerHeight+this.options.itemHeight*this.options.coverArea2Top+this.options.cssUnit}this.oneLevelContainDom=this.iosSelectLayer.el.querySelector("#oneLevelContain");
this.twoLevelContainDom=this.iosSelectLayer.el.querySelector("#twoLevelContain");this.threeLevelContainDom=this.iosSelectLayer.el.querySelector("#threeLevelContain");this.fourLevelContainDom=this.iosSelectLayer.el.querySelector("#fourLevelContain");this.fiveLevelContainDom=this.iosSelectLayer.el.querySelector("#fiveLevelContain");this.sixLevelContainDom=this.iosSelectLayer.el.querySelector("#sixLevelContain");this.oneLevelUlContainDom=this.iosSelectLayer.el.querySelector(".select-one-level");this.twoLevelUlContainDom=
this.iosSelectLayer.el.querySelector(".select-two-level");this.threeLevelUlContainDom=this.iosSelectLayer.el.querySelector(".select-three-level");this.fourLevelUlContainDom=this.iosSelectLayer.el.querySelector(".select-four-level");this.fiveLevelUlContainDom=this.iosSelectLayer.el.querySelector(".select-five-level");this.sixLevelUlContainDom=this.iosSelectLayer.el.querySelector(".select-six-level");this.iosSelectLayer.el.querySelector(".layer").style.height=this.options.itemHeight*this.options.itemShowCount+
this.options.headerHeight+this.options.cssUnit;this.oneLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;document.documentElement.classList.add("ios-select-body-class");this.scrollOne=new l("#oneLevelContain",{probeType:3,bounce:false});this.setScorllEvent(this.scrollOne,1);if(this.level>=2){this.twoLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;this.scrollTwo=new l("#twoLevelContain",{probeType:3,
bounce:false});this.setScorllEvent(this.scrollTwo,2)}if(this.level>=3){this.threeLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;this.scrollThree=new l("#threeLevelContain",{probeType:3,bounce:false});this.setScorllEvent(this.scrollThree,3)}if(this.level>=4){this.fourLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;this.scrollFour=new l("#fourLevelContain",{probeType:3,bounce:false});this.setScorllEvent(this.scrollFour,
4)}if(this.level>=5){this.fiveLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;this.scrollFive=new l("#fiveLevelContain",{probeType:3,bounce:false});this.setScorllEvent(this.scrollFive,5)}if(this.level>=6){this.sixLevelContainDom.style.height=this.options.itemHeight*this.options.itemShowCount+this.options.cssUnit;this.scrollSix=new l("#sixLevelContain",{probeType:3,bounce:false});this.setScorllEvent(this.scrollSix,6)}this.selectBtnDom=this.iosSelectLayer.el.querySelector(".sure");
this.selectBtnDom.addEventListener("click",function(){a.options.callback&&a.options.callback(a.selectOneObj,a.selectTwoObj,a.selectThreeObj,a.selectFourObj,a.selectFiveObj,a.selectSixObj)})},mapKeyByIndex:function(a){var b={index:1,levelContain:this.oneLevelContainDom,relation:this.options.oneTwoRelation};a===2?b={index:2,levelContain:this.twoLevelContainDom,relation:this.options.twoThreeRelation}:a===3?b={index:3,levelContain:this.threeLevelContainDom,relation:this.options.threeFourRelation}:a===
4?b={index:4,levelContain:this.fourLevelContainDom,relation:this.options.fourFiveRelation}:a===5?b={index:5,levelContain:this.fiveLevelContainDom,relation:this.options.fiveSixRelation}:a===6&&(b={index:6,levelContain:this.sixLevelContainDom,relation:0});return b},setScorllEvent:function(a,b){var c=this,f=c.mapKeyByIndex(b);a.on("scrollStart",function(){c.toggleClassList(f.levelContain)});a.on("scroll",function(){if(!isNaN(this.y)){var a=Math.abs(this.y/c.baseSize)/c.options.itemHeight,b=1,b=Math.round(a)+
1;c.toggleClassList(f.levelContain);c.changeClassName(f.levelContain,b)}});a.on("scrollEnd",function(){var d=Math.abs(this.y/c.baseSize)/c.options.itemHeight,e=1,g=0;if(Math.ceil(d)===Math.round(d)){g=Math.ceil(d)*c.options.itemHeight*c.baseSize;e=Math.ceil(d)+1}else{g=Math.floor(d)*c.options.itemHeight*c.baseSize;e=Math.floor(d)+1}a.scrollTo(0,-g,0);c.toggleClassList(f.levelContain);d=c.changeClassName(f.levelContain,e);e=n.attrToData(d,e);c.setSelectObj(b,e);c.level>b&&(f.relation===1&&n.isArray(c.data[b])||
n.isFunction(c.data[b]))&&c.setLevelData(b+1,c.selectOneObj.id,c.selectTwoObj.id,c.selectThreeObj.id,c.selectFourObj.id,c.selectFiveObj.id,c.selectSixObj.id)});a.on("scrollCancel",function(){var d=Math.abs(this.y/c.baseSize)/c.options.itemHeight,e=1,g=0;if(Math.ceil(d)===Math.round(d)){g=Math.ceil(d)*c.options.itemHeight*c.baseSize;e=Math.ceil(d)+1}else{g=Math.floor(d)*c.options.itemHeight*c.baseSize;e=Math.floor(d)+1}a.scrollTo(0,-g,0);c.toggleClassList(f.levelContain);d=c.changeClassName(f.levelContain,
e);e=n.attrToData(d,e);c.setSelectObj(b,e);c.level>b&&(f.relation===1&&n.isArray(c.data[b])||n.isFunction(c.data[b]))&&c.setLevelData(b+1,c.selectOneObj.id,c.selectTwoObj.id,c.selectThreeObj.id,c.selectFourObj.id,c.selectFiveObj.id,c.selectSixObj.id)})},loadingShow:function(){this.options.showLoading&&(this.iosSelectLoadingBoxDom.style.display="block")},loadingHide:function(){this.iosSelectLoadingBoxDom.style.display="none"},mapRenderByIndex:function(a){var b={index:1,relation:0,levelUlContainDom:this.oneLevelUlContainDom,
scrollInstance:this.scrollOne,levelContainDom:this.oneLevelContainDom};a===2?b={index:2,relation:this.options.oneTwoRelation,levelUlContainDom:this.twoLevelUlContainDom,scrollInstance:this.scrollTwo,levelContainDom:this.twoLevelContainDom}:a===3?b={index:3,relation:this.options.twoThreeRelation,levelUlContainDom:this.threeLevelUlContainDom,scrollInstance:this.scrollThree,levelContainDom:this.threeLevelContainDom}:a===4?b={index:4,relation:this.options.threeFourRelation,levelUlContainDom:this.fourLevelUlContainDom,
scrollInstance:this.scrollFour,levelContainDom:this.fourLevelContainDom}:a===5?b={index:5,relation:this.options.fourFiveRelation,levelUlContainDom:this.fiveLevelUlContainDom,scrollInstance:this.scrollFive,levelContainDom:this.fiveLevelContainDom}:a===6&&(b={index:6,relation:this.options.fiveSixRelation,levelUlContainDom:this.sixLevelUlContainDom,scrollInstance:this.scrollSix,levelContainDom:this.sixLevelContainDom});return b},getLevelData:function(a,b,c,d,g,e){var h=[],k=this.mapRenderByIndex(a);
if(a===1)h=this.data[0];else if(k.relation===1){var j=arguments[a-1];this.data[a-1].forEach(function(a){a.parentId==j&&h.push(a)})}else h=this.data[a-1];return h},setLevelData:function(a,b,c,d,g,e,h){if(n.isArray(this.data[a-1])){var k=this.getLevelData(a,b,c,d,g);this.renderLevel(a,b,c,d,g,e,h,k)}else if(n.isFunction(this.data[a-1])){this.loadingShow();this.data[a-1].apply(this,[b,c,d,g,e].slice(0,a-1).concat(function(k){this.loadingHide();this.renderLevel(a,b,c,d,g,e,h,k)}.bind(this)))}else throw Error("data format error");
},renderLevel:function(a,b,c,d,g,e,h,k){var j=0,p=arguments[a];k.some(function(a){return a.id==p})||(p=k[0].id);var m="",o=this.options.itemHeight+this.options.cssUnit,m=m+this.getWhiteItem();k.forEach(function(a,b){if(a.id==p){m=m+('<li style="height: '+o+"; line-height: "+o+';"'+n.attrToHtml(a)+' class="at">'+a.value+"</li>");j=b+1}else m=m+('<li style="height: '+o+"; line-height: "+o+';"'+n.attrToHtml(a)+">"+a.value+"</li>")});var m=m+this.getWhiteItem(),l=this.mapRenderByIndex(a);l.levelUlContainDom.innerHTML=
m;l.scrollInstance.refresh();l.scrollInstance.scrollToElement(":nth-child("+j+")",0);l=this.changeClassName(l.levelContainDom,j);l=n.attrToData(l,j);this.setSelectObj(a,l);this.level>a&&this.setLevelData(a+1,this.selectOneObj.id,this.selectTwoObj.id,this.selectThreeObj.id,this.selectFourObj.id,this.selectFiveObj.id,this.selectSixObj.id)},setSelectObj:function(a,b){if(a===1)this.selectOneObj=b;else if(a===2)this.selectTwoObj=b;else if(a===3)this.selectThreeObj=b;else if(a===4)this.selectFourObj=b;
else if(a===5)this.selectFiveObj=b;else if(a===6)this.selectSixObj=b},getWhiteItem:function(){var a;a=this.options.itemHeight+this.options.cssUnit;var b='<li style="height: '+a+"; line-height: "+a+'"></li>';a=""+b;this.options.itemShowCount>3&&(a=a+b);this.options.itemShowCount>5&&(a=a+b);this.options.itemShowCount>7&&(a=a+b);return a},changeClassName:function(a,b){var c;if(this.options.itemShowCount===3){c=a.querySelector("li:nth-child("+(b+1)+")");c.classList.add("at")}else if(this.options.itemShowCount===
5){c=a.querySelector("li:nth-child("+(b+2)+")");c.classList.add("at");a.querySelector("li:nth-child("+(b+1)+")").classList.add("side1");a.querySelector("li:nth-child("+(b+3)+")").classList.add("side1")}else if(this.options.itemShowCount===7){c=a.querySelector("li:nth-child("+(b+3)+")");c.classList.add("at");a.querySelector("li:nth-child("+(b+2)+")").classList.add("side1");a.querySelector("li:nth-child("+(b+1)+")").classList.add("side2");a.querySelector("li:nth-child("+(b+4)+")").classList.add("side1");
a.querySelector("li:nth-child("+(b+5)+")").classList.add("side2")}else if(this.options.itemShowCount===9){c=a.querySelector("li:nth-child("+(b+4)+")");c.classList.add("at");a.querySelector("li:nth-child("+(b+3)+")").classList.add("side1");a.querySelector("li:nth-child("+(b+2)+")").classList.add("side2");a.querySelector("li:nth-child("+(b+5)+")").classList.add("side1");a.querySelector("li:nth-child("+(b+6)+")").classList.add("side2")}return c},setBase:function(){if(this.options.cssUnit==="rem"){var a=
window.getComputedStyle(document.documentElement,null).fontSize;try{this.baseSize=/\d+(?:\.\d+)?/.exec(a)[0]}catch(b){this.baseSize=1}}else this.baseSize=1},toggleClassList:function(a){Array.prototype.slice.call(a.querySelectorAll("li")).forEach(function(a){a.classList.contains("at")?a.classList.remove("at"):a.classList.contains("side1")?a.classList.remove("side1"):a.classList.contains("side2")&&a.classList.remove("side2")})}};"undefined"!=typeof module&&module.exports?module.exports=s:"function"==
typeof define&&define.amd?define(function(){return s}):window.IosSelect=s})();