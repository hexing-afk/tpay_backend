
console.log("common.js");


function setCookie(name, value) {
    var Days = 30;
    var exp = new Date();
    exp.setTime(exp.getTime() + Days*24*60*60*1000);
    document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();
}



// 倒计时
// id：倒计时显示的元素id
// maxtime: 倒计时时间
function countDownNew(id, maxtime) {
    var element = $("#"+id)
    var interval = setInterval(function () {
        if (maxtime >= 0) {
            minutes = Math.floor(maxtime / 60);
            seconds = Math.floor(maxtime % 60);

            if (seconds < 10) {
                seconds = "0" + seconds
            }
            if (minutes >= 10) {
                msg = minutes + ":" + seconds;
            } else {
                msg = "0" + minutes + ":" + seconds;
            }

            element.text(msg);
            --maxtime;
        } else {
            element.text("00:00");
            clearInterval(interval);
        }
    },1000);
}

// 复制到到粘贴板
// el：当前元素
// str: 粘贴的内容
// successTip: 成功后的提示信息
function copyToClipboard(el,str, successTip) {
    console.log("copyToClipboard:",str);

    const input = document.createElement('input');

    // 防止安卓中拉起键盘
    input.setAttribute('readonly', 'readonly');

    // 设置要粘贴的内容
    input.setAttribute('value', str);
    //input.setAttribute('hidden', true);

    el.appendChild(input);

    // 让input获得焦点 否则可能会复制失败
    input.focus();

    // input.select() 在 ios 中并没有选中所有内容,
    //使用 inputSelectionRange 代替
    input.setSelectionRange(0, 9999);

    if (document.execCommand('copy')) {
        document.execCommand('copy');
        alert(successTip);
    }
    el.removeChild(input);
}