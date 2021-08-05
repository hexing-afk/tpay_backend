define("phone.app", ["up.m.mvc", "up.m.encrypt", "phone.i18", "ios.select"], function (j, t, f, v) {
    window.appType = "2"; var c = {}; c.requestBaseURL = window.http_proto + location.host + "/mobile/"; c.baseURL = c.requestBaseURL; c.sdkRequestBaseURL = window.http_proto + location.host + "/mobile/sdk/"; -1 != window.location.href.indexOf("authPayIndex.action") && (c.requestBaseURL = window.http_proto + location.host + "/mobile/authPay/"); -1 != window.location.href.indexOf("bindPayIndex.action") && (c.requestBaseURL = window.http_proto + location.host +
        "/mobile/bindPay/"); -1 != window.location.href.indexOf("transfer.action") && (c.requestBaseURL = window.http_proto + location.host + "/mobile/jf/"); -1 != window.location.href.indexOf("verify.action") && (c.requestBaseURL = window.http_proto + location.host + "/mobile/verify/"); -1 != window.location.href.indexOf("feeActivateIndex.action") && (c.requestBaseURL = window.http_proto + location.host + "/mobile/feeActivate/"); -1 != window.location.href.indexOf("securePayIndex.action") && (c.requestBaseURL = window.http_proto + location.host +
            "/mobile/securePay/"); c.constants = {
                langStrArray: ["zh_CN", "en_US"], transNumber: "transNumber", orderTips: { foldText: f.collapse, unfoldText: f.expand }, requestURL: {
                    getOrderInfoURL: c.requestBaseURL + "init.action", getBaseOrderInfoURL: c.requestBaseURL + "getOrderInfo.action", cardValidateURL: c.requestBaseURL + "cardValidate.action", getCardInfoURL: c.requestBaseURL + "getCardInfo.action", getCouponInfoURL: "/getCouponInfo.action", getTopPointInfoURL: "/getTopPointInfo.action", getTopPointProcessingURL: c.requestBaseURL + "getTopPointProcessing.action",
                    getRestrictPayDisplayCardInfo: c.requestBaseURL + "getRestrictPayDisplayCardInfo.action", getCrbePayCardInfoURL: c.requestBaseURL + "getCrbePayCardInfo.action", sendSMSURL: c.requestBaseURL + "sendSMS.action", sendSMSProcessing: c.requestBaseURL + "sendSMSProcessing.action", cardPayURL: c.requestBaseURL + "cardPay.action", tokenCardPayURL: c.requestBaseURL + "tokenPay.action", preAuthInitURL: c.requestBaseURL + "preauth/init.action", preAuthPayURL: c.requestBaseURL + "preauth/pay.action", preAuthQueryURL: c.requestBaseURL + "preauth/payProcessing.action",
                    preAuthResultURL: c.requestBaseURL + "preauth/result.action", fastCardPayURL: c.requestBaseURL + "proPay.action", getFastCardInfoURL: c.requestBaseURL + "getProCardInfo.action", fastCardPayProcessingURL: c.requestBaseURL + "proPayProcessing.action", cardPayProcessingURL: c.requestBaseURL + "cardPayProcessing.action", qrQueryOrder: c.requestBaseURL + "queryOrder.action", cardPayResultURL: c.requestBaseURL + "cardPayResult.action", loginURL: c.requestBaseURL + "login.action", checkLoginName: c.requestBaseURL + "checkLoginName.action",
                    checkcodeURL: c.requestBaseURL + "checkcode.action", loginoutURL: c.requestBaseURL + "logout.action", getCardListURL: c.requestBaseURL + "getCardList.action", getUserInfoURL: c.requestBaseURL + "getUserInfo.action", bankOpenStatusQuery: c.requestBaseURL + "bankOpenStatusQuery.action", bingoPromotion: c.requestBaseURL + "bingoPromotion.action", getPromotionInfo: "/getPromotionInfo.action", getInstalmentInfoURL: c.requestBaseURL + "getInstalmentInfo.action", reopenCard: c.requestBaseURL + "reopenCard.action", sdkInitInfoURL: c.requestBaseURL +
                        "sdkInit.action", sdkInitInfoURL_old: c.sdkRequestBaseURL + "sdkInit.action", sdkCardPayResultURL: c.requestBaseURL + "sdkCardPayResult.action", sdkSaveTerminalURL: c.requestBaseURL + "sdkSaveTerminal.action", checkShuMeiURL: c.requestBaseURL + "checkShuMei.action", getSmsInfoURL: c.requestBaseURL + "getSmsInfo.action"
                }, staticURL: window.mobileStaticUrl + "/resources/upop_m", analysisURL: window.http_proto + "upa.unionpaysecure.com/sampler/ubass.js?sysId=uniform_br_mob&jslib=zepto", upadmURL: "https://acpstatic.95516.com/gw/mobile/resources/upop_m/js/up/up.m.adm.js?version=2015032408",
                bankListURL: window.staticCmsUrl + "/page/gw_mobile/bankList/", encryptpd: window.http_proto + "", credentialTypeList: [{ key: "01", value: f.identityCard }, { key: "02", value: f.militaryID }, { key: "03", value: f.passport }, { key: "04", value: f.reentryPermit }, { key: "05", value: f.mTPs }, { key: "06", value: f.policeID }, { key: "07", value: f.soldiersCard }, { key: "09", value: f.foreignPassport }, { key: "12", value: f.HkMacaoIdCard }, { key: "13", value: f.TaiwanIdCard }, { key: "99", value: f.other }], areaCodeList: [{ code: "86", name: f.area_china }, { code: "852", name: f.area_hk },
                { code: "853", name: f.area_mc }, { code: "92", name: f.area_pakistan }], IDArea: {
                    11: "\u5317\u4eac", 12: "\u5929\u6d25", 13: "\u6cb3\u5317", 14: "\u5c71\u897f", 15: "\u5185\u8499\u53e4", 21: "\u8fbd\u5b81", 22: "\u5409\u6797", 23: "\u9ed1\u9f99\u6c5f", 31: "\u4e0a\u6d77", 32: "\u6c5f\u82cf", 33: "\u6d59\u6c5f", 34: "\u5b89\u5fbd", 35: "\u798f\u5efa", 36: "\u6c5f\u897f", 37: "\u5c71\u4e1c", 41: "\u6cb3\u5357", 42: "\u6e56\u5317", 43: "\u6e56\u5357", 44: "\u5e7f\u4e1c", 45: "\u5e7f\u897f", 46: "\u6d77\u5357", 50: "\u91cd\u5e86", 51: "\u56db\u5ddd", 52: "\u8d35\u5dde",
                    53: "\u4e91\u5357", 54: "\u897f\u85cf", 61: "\u9655\u897f", 62: "\u7518\u8083", 63: "\u9752\u6d77", 64: "\u5b81\u590f", 65: "\u65b0\u7586", 71: "\u53f0\u6e7e", 81: "\u9999\u6e2f", 82: "\u6fb3\u95e8", 91: "\u56fd\u5916"
                }
            }; var q = {
                userInfoHTML: '<div class="new-text-input welcome-wrapper">                        <p class="welcome">${i18.welcome},<span class="user-name">${userInfo.username}</span>&nbsp;&nbsp;<a href="" id="logout">${i18.quit}</a></p>                      </div>', areaCodeListHTML: '<div class="select-wrapper select-wrapper-areaCode inner-label-input item-areaCode" style="padding:0px 18px;margin-top:0px;border-top:none;border-bottom:none;" data-ref="areaCode">                                    <div class="input-wrapper">                                    <label class="inner-label">${i18.areaCodeNew}                                    </label>                                    <span class="up-input" id="areaCodeValue">${firstAreaCode.code} ${firstAreaCode.name}</span>                                <input type="hidden" id="areaCodeKey" value="${firstAreaCode.code}" />                                    <i class="select-area select-area-areaCode"></i>                                    <div class="dropdown-menu dropdown-menu-areaCode hide">                                    </div>                                    </div>                                </div>',
                credentialTypeHTML_New: '<div class="inner-label-input item-credentialType " data-ref="credentialType">                                            <div class="input-wrapper">                                            <span class="inner-label ">${i18.credentialType}</span>                                            <span class="up-input credentialType" id="credentialType" disabled="disabled">${firstCertType.value}</span>                                            <input type="hidden" id="firstCertTypeKey" value="${firstCertType.key}" />                                            <i class="op arrow_right" id="credentialTypeSelector"><img src="${staticURL}/zh_CN/images/phone/arrow_right.png"></i>                                            </div>                                        </div>',
                promotionHTML: '     <div class="new-select-wrapper">                                <span class="promotion-title">${i18.promotionTitle}</span>                                <div class="promotion-content" id="select-promotion" data-discountId="${firstPromotion.discountId}" data-discountSk="${firstPromotion.discountSk}">                                    <div id="select-promotion-activityNm" class="common-ellipsis select-promotion-activityNm  inline">${firstPromotion.activityNm}</div>                                    <div id="select-promotion-amt"  class="inline" style="font-weight: bold;">${firstPromotion.promotionAmt}</div>                                    <div class="select-area-new inline"><img src="${staticURL}/zh_CN/images/phone/arrow_right.png"/></div>                                </div>                            <div>',
                promotionListHTML: '<div class="promotion-list">                                <div class="common-change-shadow"></div>                                <div class="common-change-style" >                                    <div class="list-return"></div>                                    <h4>${i18.changePromotion}</h4>                                    <tpl for="promotionList">                                        <div class="common-change-data"  data-discountId="${discountId}" data-discountSk="${discountSk}"                                            data-payAmt="${payAmt}" data-activityNm="${activityNm}" data-promotionAmt="${promotionAmt}" data-available="${available}">                                            <div class="common-ellipsis promotion-activity-name">${activityNm}</div>                                            <div class="promotion-right">                                                <span style="color:#ED171F">${promotionAmt}</span>                                                <div class="list-select-icon"><img src="${staticURL}/zh_CN/images/phone/arrow_right_red.png"/></div>                                            </div>                                        </div>                                    </tpl>                                    <div class="no-use-data" id="login-for-point" style="border-bottom: 1px solid #ececec;">${i18.newLoginForPoint}                                        <div class="promotion-right" >                                            <span style="color:#ED171F">${i18.point}</span>                                            <div class="list-select-icon "><img src="${staticURL}/zh_CN/images/phone/arrow_right.png"/></div>                                         </div>                                    </div>                                    <div class="no-use-data" id="no-promotion">${i18.noPromotion}</div>                                </div>                            </div>',
                instalmentInfoHTML: '<div style="clear: both"></div>                            <div class="echo_bank" id="fq_close" >                            <div class="k-bank-list"></div>                            <div class="instalment-staging" >                            <div class="staging_title">                                <div id="instalment_cancel_img"></div>                                <div class="instalmentSelect_title">                                    ${i18.instalment_select}                                </div>                                <div id="confirm_title" class="up-btn confirm_title">                                     ${i18.confirmBtn}                                </div>                            </div>                                <div id="tab_list_instalment" class="instalment-tab_list" style="display: block">                                    <ul class="instalment-mj_ul" id="instalment_list_ul">                                        ${instalmentItemHTML}                                        <div style="clear:both"></div>                                    </ul>                                </div>                                 </div>                                </div>',
                instalmentRulesHTML: '  <div class="cvn2_Tips_Shadow"></div>                        <div class="cvn2_Tips instalment_Tips">                            <div class="instalment_Tips_Title">${i18.tips_Instalment}</div>                            <div class="instalment_Tips_SubTitle">${i18.title_repayment}</div>                            <div class="instalment_Tips_Text">${i18.text_repayment}</div>                            <div class="instalment_Tips_SubTitle">${i18.title_cancle_instalment}</div>                            <div class="instalment_Tips_Text">${i18.text_cancel_instalment}</div>                            <div class="pop-btn confirm">${i18.tips_close}</div>                        </div>',
                newInstalmentHtml: '<tpl if="${supportInstalment} == true">                            <div class="instalment-content">                                <input type="hidden" id="instalment-id"/>                                    <div id="fq_count" style="background-color: #FFFFFF">                                        <div class="instalment_count">${i18.installment_new}                                            <img id="instalment-pop" src="${staticURL}/zh_CN/images/phone/infomation.png">                                            <div id="fqkg" class="instalment_button discount_button_close"></div>                                        </div>                                        <div id="instalment-item-info" class="instalment_info hide" style="color: #666;">                                            <span>\u94f6\u8054\u7ea2\u5305 \u53ef\u62b5\u626320\u5143</span>                                            <div class="jt_left"></div>                                        </div>                                        <div style="clear:both"></div>                                    </div>                                </div>                        </tpl>',
                payCardHTML: '<div class="new-select-wrapper select-wrapper-cardList">                        <span class="paycard-title">${i18.payMethod}</span>                        <div class="paycard-content" id="select-bankcard-id" data-ref-bid="${displayCardInfo.bindId}">                            <div class="inline select-bank-name" id="select-bank-name">${bankName}</div>                            <div class="select-bank-type inline" id="select-bank-type">${displayCardInfo.cardTypeDisplay}</div>                            <div class="select-bank-suffix inline" id="select-bank-suffix">${displayCardInfo.cardNumberSuffix}</div>                            <div class="select-area-new inline "><img src="${staticURL}/zh_CN/images/phone/arrow_right.png"/></div>                        </div>                        </div>',
                bankListHTML: '<div class="echo_bank">                <div class="k-bank-list"></div>                <div class="k-bank-style" >                <div class="list-return"></div>                <h4>${i18.changeBankCard}</h4>                    <div  data-cardTypeDisplay="${displayCardInfo.cardTypeDisplay}" data-cardNumber="${displayCardInfo.cardNumberDisplay}" data-bankName="${displayCardInfo.bankName}" class="k-bank-date">                        <img src="https://acpstatic.95516.com/gw/app/resources/sdk/issuer_logo/ic_bank_${displayCardInfo.bankNo}.png"  onerror="javascript:this.src=\'${staticURL}/zh_CN/images/phone/default_bank.png\';" class="bank-logo"/>                        <div class="inline select-bank-name" style="margin-left: 36px;">${bankName}</div>                        <div class="select-bank-type inline">${displayCardInfo.cardTypeDisplay}</div>                        <div class="select-bank-suffix inline">${displayCardInfo.cardNumberSuffix}</div>                        <div class="list-select-icon inline"><img src="${staticURL}/zh_CN/images/phone/arrow_right_red.png"/></div>                    </div>                    <div class="change-card" id="change-other-payment">                        <img src="${staticURL}/zh_CN/images/phone/add_bank.png"  class="bank-logo"/>                        <div class="other-card-content"  >${i18.addNewCard}</div>                        <div class="list-select-icon"><img src="${staticURL}/zh_CN/images/phone/arrow_right.png"/></div>                    </div>                </div>                </div>',
                addNewCardCancelTipsHtml: '   <div class="addNewCard_Tips_Shadow"></div>                                <div class="addNewCard_Tips">                                    <div class="addNewCard_Tips_Text">${i18.confirmBackTips}</div>                                    <div class="pop-btn confirmBack">${i18.confirmBack}</div>                                    <div class="pop-btn cancelBack">${i18.cancelBack}</div>                                </div>', cardPayCancelTipsHtml: '   <div class="addNewCard_Tips_Shadow"></div>                                <div class="addNewCard_Tips">                                    <div class="addNewCard_Tips_Text">${i18.cardpayConfirmBackTips}</div>                                    <div class="pop-btn confirmBack">${i18.cardpayConfirmBack}</div>                                    <div class="pop-btn cancelBack">${i18.cardpayCancelBack}</div>                                </div>',
                cvn2TipsHtml: '  <div class="cvn2_Tips_Shadow"></div>                        <div class="cvn2_Tips">                            <div class="cvn2_Tips_Title">${i18.cvn2TipsTitle}</div>                            <div class="cvn2_TIps_Img"><img src="${staticURL}/zh_CN/images/phone/cvn2Tips.png"></div>                            <div class="cvn2_Tips_Text">${i18.cvn2TipsText}</div>                            <div class="pop-btn confirm">${i18.tips_iknow}</div>                        </div>'
            }; c.TPL = q;
    c.cache = {}; c.cache.TPL = {}; c.constants.userInfo = {}; var b = j.$; b.fn.upbankcard = function () {
        var e = {
            setCursor: function (b, e, c) { if (b.setSelectionRange) { b.focus(); b.setSelectionRange(e, c) } else if (b.createTextRange) { range = b.createTextRange(); range.collapse(true); range.moveEnd("character", c); range.moveStart("character", e); range.select() } }, getCursorPosition: function (b) {
                var e = 0; if (document.selection) { b.focus(); e = document.selection.createRange(); e.moveStart("character", -b.value.length); e = e.text.length } else if (b.selectionStart ||
                    b.selectionStart == "0") e = b.selectionStart; return e
            }
        }; return b(this).each(function () {
            var c = b(this).val(); b(this).val(c.replace(/\s/g, "").replace(/(\d{4})(?=\d)/g, "$1 ")); b(this).bind("keypress", function (b) { var e = b.which; return e == 0 || e == 8 || (e == 46 || e >= 48 && e <= 57) || b.ctrlKey && (e == 99 || e == 97 || e == 118 || e == 120) ? true : false }); b(this).bind("keyup", function (c) {
                c = c.which; if (c == 46 || c >= 48 && c <= 57 || c >= 96 && c <= 105) {
                    var c = e.getCursorPosition(b(this)[0]), g = 0, h = b(this).val(), g = h.split(" ").length, h = h.replace(/\s/g, "").replace(/(\d{4})(?=\d)/g,
                        "$1 "), g = h.split(" ").length - g; b(this).val(h); e.setCursor(b(this)[0], c + g, c + g); return true
                } return false
            }); b(this).change(function () { var e = b(this).val().replace(/[ ]/g, "").replace(/(\d{4})(?=\d)/g, "$1 "); b(this).val(e) }); b(this).bind("paste", function () { var e = b(this); setTimeout(function () { var c = b(e).val().replace(/[ ]/g, ""); if (/[^\d]/.test(c)) { alert(f.inputNumber); b(e).val("") } else { c = b(e).val().replace(/[ ]/g, "").replace(/(\d{4})(?=\d)/g, "$1 "); b(e).val(c) } }, 100) }); b(this).bind("dragenter", function () { return false })
        })
    };
    c.util = {
        bindReturn: function () {
            var e = b(".do-return").attr("data-locked"), g = b(".do-return").attr("data-ref"); b(".do-return").click(function () {
                b("#s_language").show(); var d = location.hash.slice(1) || "!"; if (d.indexOf("!result") == -1) {
                    if (e == "true") { if ((d.indexOf("!cardPay") != -1 || d.indexOf("!foreignPay") != -1 || d.indexOf("!cardOpen") != -1 || d.indexOf("!cardOpenAndPay") != -1 || d.indexOf("!prepaidCard") != -1 || d.indexOf("!crbePay") != -1) && g) { window.location.href = g; return } } else {
                        if (d.indexOf("!smsSend") != -1) {
                            c.constants.userInfo &&
                            c.constants.userInfo.hasOwnProperty("username") ? c.util.pageRoute({ route: "fastPay" }) : c.util.pageRoute({ route: "cardIndex" }); return
                        } if (d.indexOf("!cardChange") != -1) { c.constants.userInfo && c.constants.userInfo.hasOwnProperty("username") && b(".fast").length > 0 ? c.util.pageRoute({ route: "fastPay" }) : c.util.pageRoute({ route: "cardIndex" }); return } if (d.indexOf("!cardPay") != -1 || d.indexOf("!cardOpen") != -1 || d.indexOf("!cardOpenPay") != -1 || d.indexOf("!cardOpenAndPay") != -1 || d.indexOf("!restrictPay") != -1) {
                            if (c.constants.userInfo &&
                                c.constants.userInfo.hasOwnProperty("username")) { b("#popInfo").empty().append(b(j.template.applyTpl(q.addNewCardCancelTipsHtml, b.extend({ staticURL: c.constants.staticURL, i18: f })))).show(); b(".addNewCard_Tips .confirmBack").click(function () { b("#popInfo").hide(); c.util.pageRoute({ route: "fastPay" }) }) } else {
                                    b("#popInfo").empty().append(b(j.template.applyTpl(q.cardPayCancelTipsHtml, b.extend({ staticURL: c.constants.staticURL, i18: f })))).show(); b(".addNewCard_Tips .confirmBack").click(function () {
                                        b("#popInfo").hide();
                                        c.util.pageRoute({ route: "cardIndex" })
                                    })
                            } b(".addNewCard_Tips .cancelBack").click(function () { b("#popInfo").hide() }); b("#contentPop").css({ top: "0px" }); b(".addNewCard_Tips").show(); return
                        } if (d.indexOf("!cardIndex") == -1 && d != "!") { c.constants.userInfo = {}; c.util.pageRoute({ route: "cardIndex" }); return }
                    } if (g) window.location.href = g
                }
            })
        }, getFindPwdUrl: function () { var e; return e = "https://sign.unionpay.com/pages/wap/findpwd.html?" + b.param({ sysIdStr: "YGfjepqMep092M5", service: window.location.href }) }, getRegUrl: function () {
            var e;
            return e = "https://sign.unionpay.com/pages/wap/reg.html?" + b.param({ sysIdStr: "YGfjepqMep092M5", service: window.location.href })
        }, getReopenParam: function () { var b = location.hash.slice(1) || "!"; return b.indexOf("&isReopenCard=true") != -1 ? c.util.getNameValuePair(b.substring(b.indexOf("&isReopenCard=true") + 1)) : {} }, base64: {
            enKey: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", deKey: [-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
            -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1], encode: function (b) {
                for (var c = [], d, i, h, f = 0; f + 3 <= b.length;) {
                    d = b.charCodeAt(f++); i = b.charCodeAt(f++); h = b.charCodeAt(f++); c.push(this.enKey.charAt(d >> 2), this.enKey.charAt((d << 4) + (i >> 4) & 63)); c.push(this.enKey.charAt((i << 2) + (h >> 6) & 63), this.enKey.charAt(h &
                        63))
                } if (f < b.length) { d = b.charCodeAt(f++); c.push(this.enKey.charAt(d >> 2)); if (f < b.length) { i = b.charCodeAt(f); c.push(this.enKey.charAt((d << 4) + (i >> 4) & 63)); c.push(this.enKey.charAt(i << 2 & 63), "=") } else c.push(this.enKey.charAt(d << 4 & 63), "==") } return c.join("")
            }, decode: function (b) {
                for (var c = [], d, i, h, f, k = 0, b = b.replace(/[^A-Za-z0-9\+\/]/g, ""); k + 4 <= b.length;) {
                    d = this.deKey[b.charCodeAt(k++)]; i = this.deKey[b.charCodeAt(k++)]; h = this.deKey[b.charCodeAt(k++)]; f = this.deKey[b.charCodeAt(k++)]; c.push(String.fromCharCode((d <<
                        2 & 255) + (i >> 4), (i << 4 & 255) + (h >> 2), (h << 6 & 255) + f))
                } if (k + 1 < b.length) { d = this.deKey[b.charCodeAt(k++)]; i = this.deKey[b.charCodeAt(k++)]; if (k < b.length) { h = this.deKey[b.charCodeAt(k)]; c.push(String.fromCharCode((d << 2 & 255) + (i >> 4), (i << 4 & 255) + (h >> 2))) } else c.push(String.fromCharCode((d << 2 & 255) + (i >> 4))) } return c.join("")
            }
        }, getCurrentRoute: function () { var b = "", b = location.hash.slice(1) || "!", b = b.indexOf("!") != -1 && b.indexOf("?r=") != -1 ? b.substring(b.indexOf("!") + 1, b.indexOf("?r=")) : ""; return c.util.base64.encode(b) },
        getNameValuePair: function (b) { var c = {}; b.replace(/([^&=]*)(?:\=([^&]*))?/gim, function (b, e, h) { e != "" && (c[e] = h || "") }); return c }, serialize: function (b) { var c = [], d; for (d in b) b[d] == null || b[d] == "" ? c.push(d + "=") : c.push(d + "=" + b[d]); return c.join("&") }, handleBankName: function (b) { var c = b.length; if (c > 6) var d = b.length - 2, i = b.substring(0, 3), d = b.substring(d, c), b = i + " ... " + d; return b }, pageRoute: function (b, g) {
            var d = b.route, i = "", i = g ? g : ""; c.util.hideContainerOtherMessage(); if (d == "index") {
                c.util.showOrderDetail(); window.location =
                    "#!index?r=" + Math.random() + i
            } else if (d == "cardChange") window.location = "#!cardChange?r=" + Math.random() + i; else if (d == "smsSend") window.location = "#!smsSend?r=" + Math.random() + i; else if (d == "cardIndex") { c.util.showOrderDetail(); window.location = "#!cardIndex?r=" + Math.random() + i } else d == "fastPay" ? window.location = "#!fastPay?r=" + Math.random() + i : d == "cardPay" ? window.location = "#!cardPay?r=" + Math.random() + i : d == "cardOpen" ? window.location = "#!cardOpen?r=" + Math.random() + i : d == "cardOpenPay" ? window.location = "#!cardOpenPay?r=" +
                Math.random() + i : d == "foreignPay" ? window.location = "#!foreignPay?r=" + Math.random() + i : d == "cardOpenAndPay" ? window.location = "#!cardOpenAndPay?r=" + Math.random() + i : d == "cardOpenConfirm" ? window.location = "#!cardPay?r=" + Math.random() + i : d == "prepaidCard" ? window.location = "#!prepaidCard?r=" + Math.random() + i : d == "restrictPay" ? window.location = "#!restrictPay?r=" + Math.random() + i : d == "crbePay" ? window.location = "#!crbePay?r=" + Math.random() + i : d == "query" ? (new c.util.bankOpenStatusQuery({ qn: b.qn })).start(i) : window.location = d ==
                    "SPCardOpen" ? "#!SPCardOpen?r=" + Math.random() + i : d == "SPCardOpenAndPay" ? "#!SPCardOpenAndPay?r=" + Math.random() + i : d == "SPICBCCardOpen" ? "#!SPICBCCardOpen?r=" + Math.random() + i : d == "userLogin" ? "#!userLogin?r=" + Math.random() + i : d == "tokenPay" ? "#!tokenPay?r=" + Math.random() + i : d == "result" ? "#!result?r=" + Math.random() + i : d == "erweima" ? "#!erweima?r=" + Math.random() + i : d == "errorResultPage" ? "#!errorResultPage?r=" + Math.random() + i : "#!sysError?r=" + Math.random()
        }, bankOpenStatusQuery: function (b) {
            var g = function (d) {
                var i = "", i = d ?
                    d : "", h = function (b) { window.cardNumberLocked ? c.util.hint.init(function () { c.util.hint.removePop() }, function () { c.util.hint.removePop() }, f.update, f.iknow) : c.util.hint.init(function (b) { c.util.hint.removePop(); b.find("a").removeAttr("href"); c.util.pageRoute({ route: "cardIndex" }) }, function () { c.util.hint.removePop() }, f.changeCard, f.iknow); c.util.loadingFinish(); c.util.hint.changeHintText(b); c.util.hint.showPop() }; c.util.loadingStart(); j.ajax.post(c.util.buildAjaxRequest({
                        url: c.constants.requestURL.bankOpenStatusQuery,
                        data: function () { var c = { qn: b.qn }; if (i.indexOf("isReopenCard=true") != -1) c.isReopenCard = true; return c }(), success: function (b) { b = JSON.parse(b); if (b.r != "00") h(b.m); else if (b.p.status == "InProcess") setTimeout(function () { g(window.queryFun ? window.queryFun.appendParam : null) }, 1E3); else { b.p.route == "cardOpenConfirm" ? c.util.pageRoute({ route: "cardPay" }) : c.util.pageRoute({ route: b.p.route }, i); c.util.finishLoading() } }, error: function (b) { h(b); console.log("querybank error,code:" + b) }
                    }))
            }; return {
                start: function (b) {
                    window.queryFun =
                    { appendParam: b }; g(b)
                }
            }
        }, dealWhiteMerchantCtrl: function (e) { b("#yellow-msg").find("div").remove(); if (e.merchantCtrl && e.merchantCtrl.type == "WHITE") { b("#yellow-msg").show().append(b('<div class="whiteMerchantCtrl">' + e.merchantCtrl.tips + "</div>")).show(); b(".reopenCard").length > 0 && b(".reopenCard").click(function () { c.util.reOpenCard({ parentRoute: c.util.getCurrentRoute() }) }) } }, clearRestrictPay: function () { b("#yellow-msg").find(".restrictPay").remove(); b("#yellow-msg").find("div").length == 0 && b("#yellow-msg").hide() },
        dealRestrictPay: function (e) { if (e) { b("#yellow-msg").find(".restrictPay").remove(); b("#yellow-msg").show().append(b('<div class="restrictPay">' + f.restrictHint + '<a id="restrict-pay-open" >' + f.clickOpen + "</a></div>")); b(".restrictPay a").click(function () { c.util.reOpenCard({ parentRoute: c.util.getCurrentRoute() }) }) } else { b("#yellow-msg").find(".restrictPay").remove(); b("#yellow-msg").find("div").length == 0 && b("#yellow-msg").hide() } }, loginOut: function (b) {
            b.preventDefault(); j.ajax.post(c.util.buildAjaxRequest({
                url: c.constants.requestURL.loginoutURL,
                success: function (b) { b = JSON.parse(b); b.r != "00" ? console.log("loginout error,code:" + b.m) : c.util.pageRoute({ route: "userLogin" }) }, error: function (b) { console.log("loginout error,code:" + b) }
            }))
        }, reOpenCard: function (b) {
            c.util.loadingStart(); j.ajax.post(c.util.buildAjaxRequest({
                data: {}, url: c.constants.requestURL.reopenCard, success: function (g) {
                    c.util.finishLoading(); g = JSON.parse(g); if (g.r != "00") {
                        c.util.hint.init(function () { c.util.hint.removePop() }, function () { c.util.hint.removePop() }, f.confirmBtn, f.iknow); c.util.loadingFinish();
                        c.util.hint.changeHintText(g.m); c.util.hint.showPop()
                    } else b ? c.util.pageRoute(g.p, "&isReopenCard=true&" + c.util.serialize(b)) : c.util.pageRoute(g.p)
                }, error: function () { console && console.log("card validate\u5931\u8d25 \uff01") }
            }))
        }, buildAjaxRequest: function (e) {
            if (!e.data) e.data = {}; e.data.t = c.util.getTn(); if (!window.tokenId && j.util.cookie.getcookie("m_t_key", true)) window.tokenId = j.util.cookie.getcookie("m_t_key", true); if (window.tokenId) e.data.token = window.tokenId; return b.extend({
                async: true, dataType: "text/json",
                timeoutError: function () { c.util.finishLoading(); console && console.log("ajax \u8bf7\u6c42 \u8d85\u65f6") }
            }, e)
        }, initToast: function () { document.onscroll = function () { b(window).scrollTop() != 0 ? b(".toast").addClass("screen") : b(".toast").removeClass("screen"); b("#app_download").length > 0 && b("#app_download").css("bottom", b(window).scrollTop() * -1) } }, initZzh: function () { b.getScript(c.constants.cUrl + "?tn=" + c.util.getTn()) }, trim: function (b) { return b.replace(/\s/g, "") }, luhnChk: function (b) {
            return RegExp(/^\d{13,19}$/).test(b) ?
                true : false
        }, regex: {
            cvn2: RegExp(/^\d{3}$/), mobile: RegExp(/^1\d{10}$/), expire: RegExp(/^\d{4}$/), smsCode: RegExp(/^\d{6}$/), foreignPhone: RegExp(/^\d{4,15}$/), pass: function (b, g) { if (!g) return false; if (b == "credential") { if (c.util.regex.checkIDNumber(g)) return true } else { var d = c.util.regex[b]; if (!d || d.test(g)) return true } return false }, checkIDNumber: function (e) {
                if (b("#firstCertTypeKey").val() == "01") {
                    var g = c.constants.IDArea, d = [], d = e.split(""); g[parseInt(e.substr(0, 2))] == null && c.util.toast(f.idIllegal); switch (e.length) {
                        case 15: ereg =
                            (parseInt(e.substr(6, 2)) + 1900) % 4 == 0 || (parseInt(e.substr(6, 2)) + 1900) % 100 == 0 && (parseInt(e.substr(6, 2)) + 1900) % 4 == 0 ? /^[1-9][0-9]{5}[0-9]{2}((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)(0[1-9]|[1-2][0-9]|30)|02(0[1-9]|[1-2][0-9]))[0-9]{3}$/ : /^[1-9][0-9]{5}[0-9]{2}((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)(0[1-9]|[1-2][0-9]|30)|02(0[1-9]|1[0-9]|2[0-8]))[0-9]{3}$/; if (ereg.test(e)) return true; c.util.toast(f.idBirthDayIllegal); return false; case 18: ereg = parseInt(e.substr(6,
                                4)) % 4 == 0 || parseInt(e.substr(6, 4)) % 100 == 0 && parseInt(e.substr(6, 4)) % 4 == 0 ? /^[1-9][0-9]{5}(19|20)[0-9]{2}((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)(0[1-9]|[1-2][0-9]|30)|02(0[1-9]|[1-2][0-9]))[0-9]{3}[0-9Xx]$/ : /^[1-9][0-9]{5}(19|20)[0-9]{2}((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)(0[1-9]|[1-2][0-9]|30)|02(0[1-9]|1[0-9]|2[0-8]))[0-9]{3}[0-9Xx]$/; if (ereg.test(e)) {
                                    e = (parseInt(d[0]) + parseInt(d[10])) * 7 + (parseInt(d[1]) + parseInt(d[11])) * 9 + (parseInt(d[2]) + parseInt(d[12])) *
                                    10 + (parseInt(d[3]) + parseInt(d[13])) * 5 + (parseInt(d[4]) + parseInt(d[14])) * 8 + (parseInt(d[5]) + parseInt(d[15])) * 4 + (parseInt(d[6]) + parseInt(d[16])) * 2 + parseInt(d[7]) * 1 + parseInt(d[8]) * 6 + parseInt(d[9]) * 3; e = "10X98765432".substr(e % 11, 1); if (e == d[17]) return true; c.util.toast(f.inputIdNum); return false
                                } c.util.toast(f.idBirthDayIllegal); return false; default: c.util.toast(f.inputIdNum); return false
                    }
                } else {
                    if (/^[^!$%\^&*?<>]{5,18}$/.test(e)) return true; c.util.toast(f.inputRight + b("#firstCertTypeValue").html() + f.number);
                    return false
                }
            }
        }, validateAndCollect: function (e) {
            var g = { result: "0", data: {} }; window.isInstalmentMandated && b("#fqkg").each(function (e, i) { if (b(i).hasClass("discount_button_close")) { g.result = "1"; c.util.toast(f.makeSure + f.installment + f.switchBtn + f.isOpened); return false } }); b(".form-container .inner-label-input").each(function (d, i) {
                if (b(i).hasClass("hide") || b(i).find("input").length == 0) {
                    if (b(i).attr("data-ref") == "encpassword") {
                        if (b("#password_text").length == 0 || b("#password_text").html().length != 6) {
                            g.result =
                            "1"; c.util.toast(f.validateCode); return false
                        } g.data.password = b("#encpassword").val()
                    }
                } else {
                    var h = b(i).find("input").val(), n = b(i).attr("data-ref"); if (!b(i).find("." + n).hasClass("disabled")) {
                        if (n != "mobile" && !c.util.regex.pass(n, h)) { g.result = "1"; h = b(i).find(".inner-label").text(); n != "credential" && c.util.toast(f.makeSure + h + f.isRight); return false } if (n == "mobile") if (!b("#areaCodeKey").val() || b("#areaCodeKey").val() == "86") {
                            if (!c.util.regex.pass("mobile", h)) {
                                g.result = "1"; h = b(i).find(".inner-label").text();
                                c.util.toast(f.makeSure + h + f.isRight); return false
                            }
                        } else if (!c.util.regex.pass("foreignPhone", h)) { g.result = "1"; h = b(i).find(".inner-label").text(); c.util.toast(f.makeSure + h + f.isRight); return false }
                    } if (n == "mobile" && e && e.bankPhoneNumber) if (e.bankPhoneNumber.substring(0, 3) == h.substring(0, 3) && e.bankPhoneNumber.substring(7, 11) == h.substring(7, 11)) g.result = "0"; else { g.result = "1"; h = b(i).find(".inner-label").text(); c.util.toast(f.makeSure + h + f.isRight); return false } n == "password" && (h = t.a(h, c.constants.sk)); n ==
                        "cvn2" && (g.data.enCvn2 = t.a(h, c.constants.sk)); g.data[n] = h
                }
            }); b(".agreement-wrapper4CQP").length > 0 && (g.data.isCQPAgreementChecked = b(".agreement-wrapper4CQP").hasClass("agreed")); return g
        }, handleOrderFolder: function () { b(".order-list").find(".fold").toggleClass("hide"); b(".btn-show-order").text(); b(this).hasClass("order-list-block-xl") ? b(this).removeClass().addClass("order-list-block-down") : b(this).removeClass().addClass("order-list-block-xl") }, foldOrderDetail: function () {
            if (b("#order-list-block-show").hasClass("order-list-block-down")) {
                b(".order-list").find(".fold").toggleClass("hide");
                b("#order-list-block-show").removeClass().addClass("order-list-block-xl")
            }
        }, setOrderStatus: function (c, g) { b(".order-status").text(c).addClass(g).toggleClass("hide") }, hideOrderStatus: function () { b(".order-status").addClass("hide"); b("#error-msg ,#yellow-msg").css("display", "none") }, hideOrderDetail: function () { b("#order-list-block-show").removeClass().addClass("order-list-block-xl").addClass("hide"); b(".order-list").find(".fold").addClass("hide") }, hidePaytimeoutTips: function () {
            navigator.userAgent.toLowerCase().indexOf("android") !=
            -1 ? b("#title_tips").removeClass("title_android_with_paytiemout").addClass("title_android").show() : b("#title_tips").removeClass("title_with_paytiemout").addClass("title").show(); b("#paytimeout_tips").css("display", "none")
        }, showOrderDetail: function () { b(".order-list-block-xl").removeClass("hide"); b(".order-list-block-down").removeClass("hide") }, initPromotion: function (e) {
            j.$("#new_promotion").html(j.template.applyTpl(q.promotionHTML, {
                staticURL: c.constants.staticURL, promotions: e.promotions, firstPromotion: e.firstPromotion,
                i18: f
            })).show(); if (e.firstPromotion && e.firstPromotion.available) { var g = b("#order-item-origin-amount").text(), e = e.firstPromotion.payAmt; if (g == e) { b(".order-item-origin-amount").hide(); b("#order-item-pay-amount").text(g) } else { b(".order-item-origin-amount").show(); b("#order-item-pay-amount").text(e) } }
        }, bindPromotionEvent: function (e) {
            complete(); b("#popInfo").empty().append(b(j.template.applyTpl(q.promotionListHTML, { staticURL: c.constants.staticURL, promotionList: e.promotions.promotionList, i18: f }))).show();
            b(".common-change-data").each(function () { b(this).find("img").hide(); b(this).attr("data-discountId") == b("#select-promotion").attr("data-discountId") && b(this).find("img").show(); if (b(this).attr("data-available") == "false") { b(this).find(".promotion-right").hide(); b(this).addClass("notclick"); b(this).css("color", "#d3d3d3") } }); if (c.constants.userInfo && c.constants.userInfo.hasOwnProperty("username")) b("#login-for-point").hide(); else { b("#login-for-point").show(); b("#login-for-point").click(function () { c.util.pageRoute({ route: "userLogin" }) }) } b(".list-return").click(function () { b("#popInfo").hide() });
            b(".common-change-data").click(function () {
                if (b(this).attr("data-discountId") != b("#select-promotion").attr("data-discountId")) {
                    Instalment.close(); var c = b("#order-item-origin-amount").text(), e = b(this).attr("data-payAmt"); if (c == e) { b(".order-item-origin-amount").hide(); b("#order-item-pay-amount").text(c) } else { b(".order-item-origin-amount").show(); b("#order-item-pay-amount").text(e) } b("#select-promotion-activityNm").text(b(this).attr("data-activityNm")); b("#select-promotion-amt").text(b(this).attr("data-promotionAmt"));
                    b("#select-promotion").attr("data-discountId", b(this).attr("data-discountId")); b("#select-promotion").attr("data-discountSk", b(this).attr("data-discountSk")); b(".promotion-content").css("color", "#ED171F")
                } console.log(Math.round(b("#order-item-pay-amount").text() * 100)); b("#popInfo").hide()
            }); b("#no-promotion").click(function () {
                b("#order-item-pay-amount").text(b("#order-item-origin-amount").text()); b(".order-item-origin-amount").hide(); b("#select-promotion-activityNm").text(f.noPromotion); b("#select-promotion-amt").text("");
                b(".promotion-content").attr("data-discountId", ""); b(".promotion-content").attr("data-discountSk", ""); b(".promotion-content").css("color", "#666666"); b("#popInfo").hide(); Instalment.close()
            }); b("#contentPop").css({ top: "0px" })
        }, clearPromotion: function () { b(".promotion-list").hide(); b("#order-item-pay-amount").text(b("#order-item-origin-amount").text()); b(".order-item-origin-amount").hide(); b(".promotion-content").attr("data-discountId", ""); b(".promotion-content").attr("data-discountSk", "") }, initInstalment: function (b) {
            j.$("#new_instalment").html(j.template.applyTpl(c.TPL.newInstalmentHtml,
                { staticURL: c.constants.staticURL, supportInstalment: b.supportInstalment, i18: f })).show()
        }, hideContainerOtherMessage: function () { b(".promotion-list").hide(); b(".paycard-list").hide(); b(".instalment-list").hide(); b("#order-item-pay-amount").text(b("#order-item-origin-amount").text()); b(".order-item-origin-amount").hide(); b(".promotion-content").attr("data-discountId", ""); b(".promotion-content").attr("data-discountSk", ""); Instalment.close(); b("#yellow-msg").hide(); b("#error-msg").hide() }, copyValue: function (c) {
            ls =
            window.setInterval(function () { var g = c.value.replace(/\D|\s/g, "").replace(/(\d{4})(?=\d)/g, "$1 "); c.value.length == 0 ? b(".op.del").addClass("hide") : b(".op.del").removeClass("hide"); if (g != c.value) c.value = g; document.getElementsByClassName("pan").innerText = c.value }, 200)
        }, toggleAgreement: function (c) { c.preventDefault(); if (c = b(c.target).attr("href")) { window.location.href = c; return false } b(".agreement-wrapper").toggleClass("agreed"); b(".agreement-wrapper").find("span").toggleClass("hide") }, toggleRememberCardAgreement: function (c) {
            c.preventDefault();
            if (c = b(c.target).attr("href")) { window.location.href = c; return false } b(".agreement-wrapper-rememberCard").toggleClass("agreed"); b(".agreement-wrapper-rememberCard").find("span").toggleClass("hide")
        }, log: function (b) { console.log(JSON.stringify(b)) }, credentialTypeBindEvent: function () {
            var e = b("#credentialType"); b("#credentialTypeSelector,#credentialType").click(function () {
                new v(1, [c.constants.credentialTypeList.map(function (b) { return { id: b.key, value: b.value } })], {
                    title: "", container: "#select_wrapper", itemHeight: 36,
                    headerHeight: 54, itemShowCount: 5, oneLevelId: e, sureText: f.done, closeText: f.cancel, callback: function (c) { b("#credentialType").html(c.value); b("#firstCertTypeKey").val(c.id) }
                })
            })
        }, cvn2TipsBindEvent: function () { b(".questionMark").length > 0 && b(".questionMark").click(function () { b("#popInfo").empty().append(b(j.template.applyTpl(c.TPL.cvn2TipsHtml, { staticURL: c.constants.staticURL, i18: f }))).show(); b(".cvn2_Tips .confirm").click(function () { b("#popInfo").hide() }); b("#contentPop").css({ top: "0px" }); b(".cvn2_Tips").show() }) },
        getTn: function () { var b = ""; try { b = j.util.url(window.location.href).getParam(c.constants.transNumber) } catch (g) { } if (!b) b = window.transNumber; return b }, getSign: function () { return j.util.url(window.location.href).query.sign }, hint: {
            downLoadInit: function (b) { navigator.userAgent.match(/Android/i) ? b.find("a").attr("href", "https://mpay.unionpay.com/getclient?platform=android&type=cashier") : b.find("a").attr("href", "https://mpay.unionpay.com/getclient?platform=ios&type=cashier") }, showPop: function () { b(".popup-block").removeClass("hide") },
            removePop: function () { b(".popup-block").addClass("hide") }, togglePop: function () { b(".popup-block").toggleClass("hide") }, init: function (c, g, d, i) { b(".popup-block .popup-left").off("click"); b(".popup-block .popup-right").off("click"); b(".popup-block .popup-left").click(function () { c(b(this)) }).find("a").text(d); b(".popup-block .popup-right").click(function () { g(b(this)) }).find("a").text(i) }, changeHintText: function (c) { b(".popup-block .hint").html(c) }
        }, pageInitFailed: function (e) {
            b(".order-list-wrapper .order-status").text(e).removeClass("hide");
            c.util.finishLoading()
        }, bindSms: function (e) {
            window.countDownHandle && clearInterval(window.countDownHandle); b("#sendCode").click(function () {
                var g = this, d = e(); if (d.ret == "1") { c.util.toast(f.validateMobile); return false } if (d.ret == "2") { c.util.toast(f.validateZtyt); return false } var i = {}; i.data = d.params; b(g).off("click"); c.flow.sms(i, function () { i.data.bindId && i.data.bindId != b("#select-bankcard-id").attr("data-ref-bid") || c.util.toast(f.sendSuccess); c.util.countDown(b(g), e, i) }, function (d) {
                    c.util.bindSms(e);
                    if (d && d.indexOf("reopenCard57") != -1) { c.util.hint.changeHintText(d); c.util.hint.showPop(); b(".reopenCard57").length > 0 && b(".reopenCard57").click(function () { c.util.reOpenCard({}); c.util.hint.removePop() }) } else c.util.toast(d)
                }, function () { c.util.toast(f.sendFailed); c.util.bindSms(e) })
            }); return this
        }, countDown: function (c, g, d) {
            var i = this, h = 60; b("#sendCode").length > 0 && b("#sendCode").attr("disabled", true).removeClass("sms").addClass("smsDisable") && b("#sendCode").off("click"); b("#sendCode").text(h + "s"); window.countDownHandle &&
                clearInterval(window.countDownHandle); window.countDownHandle = setInterval(function () { h = h - 1; b("#sendCode").text(h + "s"); d.data.bindId && d.data.bindId != b("#select-bankcard-id").attr("data-ref-bid") && (h = 0); if (h == 0) { b("#sendCode").text(f.getSmsCode); b("#sendCode").length > 0 && b("#sendCode").removeAttr("disabled").removeClass("smsDisable").addClass("sms"); clearInterval(window.countDownHandle); h = 60; i.bindSms(g) } }, 1E3)
        }, toast: function (c, g) {
            b(".toast").text(c).toggleClass("hide"); g ? b(".toast").addClass(g) : b(".toast").removeClass("instalment_toast");
            setTimeout(function () { b(".toast").toggleClass("hide") }, 2E3)
        }, checkEnable: function (c, g) {
            setInterval(function () {
                var d = true; c.each(function (c, e) { if (!b(e).hasClass("hide")) { var g = b(e).find("input"); if (g.length != 0) { var f = g.hasClass("disabled"); g.val().length == 0 ? b(e).find(".op.del").addClass("hide") : f || b(e).find(".op.del").removeClass("hide"); g = b(e).find("input").val(); d = d && g } } }); if (d) {
                    g.removeClass("disable"); if (b(".header_red").length > 0) {
                        g.css("background-image", "-webkit-linear-gradient(red, red)"); g.css("border",
                            "1px solid red"); g.css("background-color", "red")
                    }
                } else g.addClass("disable")
            }, 100)
        }, loadingStart: function () { b(".loading-block").removeClass("hide") }, loadingFinish: function () { b(".container").removeClass("hide"); b(".loading-block").addClass("hide") }, finishLoading: function () { b(".container").removeClass("hide"); b(".loading-block").addClass("hide") }
    }; c.Msg = {
        cmd: { init: "init", order: "order", cardbin: "cardbin", rules: "rules" }, build: function (e, g) { var d = { t: c.util.getTn(), params: {} }; g && b.extend(d.params, g); return d },
        setParams: function (c, g, d) { var i = {}; i[g] = d; b.extend(c.params, i) }, putAll: function (c, g) { b.extend(c.params, g) }, getParams: function (b, c) { return b.params[c] }, parse: function (b) { return JSON.parse(b) }
    }; c.flow = {
        querySendSMSResult: function (e, g, d, i) {
            var h = b.extend({ cnt: 1, interval: 2E3, MAX_CNT: 30, fireExceedMaxCntEvent: function () { i() } }, { cnt: e.cnt }); j.ajax.post(c.util.buildAjaxRequest({
                url: c.constants.requestURL.sendSMSProcessing, data: e.data, success: function (f) {
                    f = JSON.parse(f); if (f.r != "00") d(f.m); else {
                        var k = f.p.status;
                        e.data.bindId && e.data.bindId != b("#select-bankcard-id").attr("data-ref-bid") ? g() : k == "InProcess" ? setTimeout(function () { h.cnt++; console.log(h.cnt); if (h.cnt == h.MAX_CNT) h.fireExceedMaxCntEvent(); else { e.cnt = h.cnt; c.flow.querySendSMSResult(e, g, d, i) } }, h.interval) : k == "Succeed" ? g(f) : d(f.m)
                    }
                }, error: i
            }))
        }, queryResult: function (e, g, d, i) {
            var h = b.extend({ cnt: 1, interval: 2E3, MAX_CNT: 30, fireExceedMaxCntEvent: function () { } }, i || {}), f = c.constants.requestURL.cardPayProcessingURL; if (i.payType == "proPay") f = c.constants.requestURL.fastCardPayProcessingURL;
            if (i.payType == "preAuth") f = c.constants.requestURL.preAuthQueryURL; j.ajax.post(c.util.buildAjaxRequest({ url: f, data: e.data, success: function (b) { b = JSON.parse(b); b.r != "00" ? d(b.m) : b.p.status == "InProcess" ? setTimeout(function () { h.cnt++; console.log(h.cnt); h.cnt == h.MAX_CNT ? b && b.p && b.p.isBillPayment ? g(b) : i.fireExceedMaxCntEvent() : c.flow.queryResult(e, g, d, h) }, h.interval) : g(b) }, error: function () { setTimeout(function () { c.util.pageRoute({ route: "result" }) }, h.MAX_CNT * 1E3) } }))
        }, sms: function (e, g, d, i) {
            j.ajax.post(c.util.buildAjaxRequest({
                url: c.constants.requestURL.sendSMSURL,
                data: e.data, success: function (h) { h = JSON.parse(h); if (h.r != "00") d(h.m); else if (e.data.bindId && e.data.bindId != b("#select-bankcard-id").attr("data-ref-bid")) g(); else if (h.p.proceed) { var f = { data: {} }; f.data.qn = h.p.qn; if (e.data.bindId) f.data.bindId = e.data.bindId; c.flow.querySendSMSResult(f, g, d, i) } else g() }, error: i
            }))
        }, trans: function (b, g, d, i) {
            if (b.data.name) b.data.name = encodeURI(b.data.name); j.ajax.post(c.util.buildAjaxRequest({
                url: b.url, data: b.data, success: function (b) {
                    b = JSON.parse(b); b.r != "00" ? d(b.m) : c.flow.queryResult({ data: {} },
                        g, d, i)
                }, error: d
            }))
        }, send: function (b, g, d) { j.ajax.post(c.util.buildAjaxRequest({ url: c.constants.requestURL.getFastCardInfoURL, data: b.data, success: function (b) { b = JSON.parse(b); b.r != "00" && b.r != "G1" ? d(b.m) : g(b) }, error: d })) }
    }; var r = function (b) { var c = 0, d; for (d in o) b != d && (c = c + o[d]); return c }, o = null, s; window.PayDiscount = {
        init: function (b, c) { s = b; o = c || {}; return r() < s }, set: function (b, c) { c = parseInt(c, 10); o[b] = c; return r() < s }, remove: function (b) { if (o && o[b]) { delete o[b]; return true } return false }, get: function (b) { return o[b] },
        sum: r, left: function (b) { return s - r(b) }, removeAll: function () { for (var b in o) delete o[b] }, showAmount: function () { }, TYPE_UPOINT: "upoint", TYPE_TOPPOINT: "topPoint", TYPE_DISCOUNT: "discount", TYPE_COUPON: "coupon"
    }; var p = {}; window.UPointRelative = { TYPE_UPOINT: "uPoint", TYPE_TOPPOINT: "topPoint", TYPE_DISCOUNT: "discount", TYPE_COUPON: "coupon", use: function (b) { p[b] = 1 }, unuse: function (b) { if (p[b]) { delete p[b]; return true } return false }, unuseAll: function () { p = {} }, hasUsed: function (c) { return c ? !!p[c] && p[c] == 1 : !b.isEmptyObject(p) } };
    window.UPBaseWidget = function () { this.tplParse = function (b, c) { return b.replace(/\{([^}]+)\}/ig, function (b, e) { return c[e] != null && c[e] != void 0 ? c[e] : b }) } }; var l = function () { window.UPBaseWidget.apply(this, arguments); return this._init.apply(this, arguments) }; l.prototype = {
        _init: function (e) {
            this.settings = b.extend({}, e || {}); if (!this.settings.supportPromotion) return this; if (this.settings.promotions) { var f = this.settings; if (f.promotions.checked) { c.util.initPromotion(f); b(".promotion-content").click(function () { c.util.bindPromotionEvent(f) }) } } else b.getJSON(window.promotionUrl +
                "?" + function (b) { var c = [], e; for (e in b) b[e] && c.push(encodeURIComponent(e) + "=" + encodeURIComponent(b[e])); c.push(("v=" + Math.random()).replace(".", "")); return c.join("&") }({ bindId: this.settings.bindId, payType: this.settings.payType, mobile: this.settings.mobile, t: c.util.getTn() }) + "&callback=?", function (e) {
                    try {
                        c.util.loadingFinish(); if (e.r == "00") { var f = e.p; if (f.promotions.checked) { c.util.initPromotion(f); b(".promotion-content").click(function () { c.util.bindPromotionEvent(f) }) } else c.util.clearPromotion() } else if (e.r ==
                            "S1") window.location = "#!sysError?r=" + Math.random()
                    } catch (g) { c.util.loadingFinish() } window.isForceInstalment && (b("#fqkg").hasClass("discount_button_close") ? b("#fqkg").trigger("click") : window.Instalment.instance.refresh())
                }); return this
        }, _clear: function () { }
    }; l.instance = null; l.show = function (b) { l.options = b; l.instance = new l(b); return l.instance }; l.getCurrentDiscountItem = function (b) { return !l.instance ? null : l.instance._getCurrentDiscountItem(b) }; l.isSupport = function () {
        var b = false, b = location.hash.slice(1) ||
            "!", b = b.substring(b.indexOf("!") + 1, b.indexOf("?r=")); return (b = b.indexOf("cardPay") != -1 || b.indexOf("restrictPay") != -1 ? true : b.indexOf("cardOpen") != -1 ? b.indexOf("cardOpenAndPay") != -1 ? true : b.indexOf("cardOpenPay") != -1 ? true : false : b.indexOf("prepaidCard") != -1 ? true : b.indexOf("fastPay") != -1 ? true : b.indexOf("foreignPay") != -1 ? true : b.indexOf("crbePay") != -1 ? true : false) && !!c.constants.supportPromotion
    }; l.getPostData = function () {
        var c = b("#select-promotion").attr("data-discountId"), f = b("#select-promotion").attr("data-discountSk"),
        d = {}; if (c && c != "") if (f && f != "") { d.discountId = c; d.discountSK = f } else d.discountId = c; return d
    }; window.DisCount = l; var m = function () { window.UPBaseWidget.apply(this, arguments); return this._init.apply(this, arguments) }; m.prototype = {
        _init: function (e) { this.settings = b.extend({}, e || {}); c.util.initInstalment(this.settings); m.close(); this._bindEvent(); window.isForceInstalment && b("#fqkg").click(); return this }, refresh: function () {
            var e = this, g = b("#order-item-pay-amount").text() * 100; c.constants.supportPromotion && b("#yhkg").hasClass("discount_button_open") &&
                (g = window.PayDiscount.left()); g = Math.round(g); j.ajax.post(c.util.buildAjaxRequest({
                    data: { realAmount: g }, url: c.constants.requestURL.getInstalmentInfoURL, success: function (d) {
                        d = JSON.parse(d); c.util.loadingFinish(); if (d.p.instalmentOptions != "unsupport" && d.p.instalmentOptions != "nodata") {
                            e.settings.instalmentOptions = d.p.instalmentOptions; var d = e.settings.instalmentOptions[0], g; g = d.feeType; var h = f.instalment_rule + g; g == "average" ? h = j.template.applyTpl("{periods} ${i18.installment_term}|${i18.installment_fee}:{feeDisplay} ${i18.rmb}|${i18.installment_feeRate}:{feeRateDisplay}%|${i18.installment_average}:{averageDisplay} ${i18.rmb},${i18.installment_last} {lastDisplay} ${i18.rmb}",
                                b.extend(d, { i18: f, orderCurrency: c.constants.orderInfo.orderCurrency })) : g == "allInfirst" && (h = j.template.applyTpl("{periods} ${i18.installment_term}|${i18.installment_fee}:{feeDisplay} ${i18.rmb}|${i18.installment_feeRate}:{feeRateDisplay}%|${i18.installment_first}:{firstDisplay} ${i18.rmb},${i18.installment_lastAverage} {averageDisplay} ${i18.rmb}", b.extend(d, { i18: f, orderCurrency: c.constants.orderInfo.orderCurrency }))); g = h; d = d.periods; b("#instalment-item-info").attr("periods", d); window.selectInstalmentPeriods =
                                    d; window.selectInstalmentContent = g; try { var n = g.split("|"), k = n[0] + " ", l = n[3].split(",")[0]; b("#instalment-item-info").html("<span>" + k + l + '</span><div class="jt_left"></div>').show() } catch (u) { b("#instalment-item-info").html("<span>" + g + '</span><div class="jt_left"></div>').show() }
                        } else if (d.p.instalmentOptions == "nodata") { e.settings.instalmentOptions = ""; b("#instalment-item-info").attr("periods", null); m.close(); c.util.toast(f.installmentOutOfAmount, "instalment_toast") } else {
                            e.settings.instalmentOptions =
                            ""; b("#instalment-item-info").attr("periods", null); c.util.toast(f.notSupportInstalment, "instalment_toast"); m.close()
                        }
                    }, error: function () { c.util.loadingFinish(); console && console.log("card validate\u5931\u8d25 \uff01") }
                }))
        }, _bindEvent: function () {
            var e = this, g = function (d) {
                var e = d.feeType, g = f.instalment_rule + e; e == "average" ? g = j.template.applyTpl("{periods} ${i18.installment_term}|${i18.installment_fee}:{feeDisplay} ${i18.rmb}|${i18.installment_feeRate}:{feeRateDisplay}%|${i18.installment_average}:{averageDisplay} ${i18.rmb},${i18.installment_last} {lastDisplay} ${i18.rmb}",
                    b.extend(d, { i18: f, orderCurrency: c.constants.orderInfo.orderCurrency })) : e == "allInfirst" && (g = j.template.applyTpl("{periods} ${i18.installment_term}|${i18.installment_fee}:{feeDisplay} ${i18.rmb}|${i18.installment_feeRate}:{feeRateDisplay}%|${i18.installment_first}:{firstDisplay} ${i18.rmb},${i18.installment_lastAverage} {averageDisplay} ${i18.rmb}", b.extend(d, { i18: f, orderCurrency: c.constants.orderInfo.orderCurrency }))); return g
            }; b("#instalment-pop").unbind("click"); b("#instalment-pop").click(function (d) {
                d.preventDefault();
                b("#popInfo").empty().append(b(j.template.applyTpl(q.instalmentRulesHTML, b.extend({ staticURL: c.constants.staticURL, i18: f })))).show(); b("#popInfo").show(); b(".instalment_Tips .confirm").click(function () { b("#popInfo").hide() }); b("#contentPop").css({ top: "0px" }); b(".instalment_Tips").show()
            }); b("#fqkg").unbind("click"); b("#fqkg").click(function () {
                var d = document.getElementById("fqkg"); if (b("#fqkg").hasClass("discount_button_close")) {
                    window.selectInstalmentPeriods = ""; window.selectInstalmentContent = ""; var i =
                        b("#order-item-pay-amount").text() * 100; c.constants.supportPromotion && b("#yhkg").hasClass("discount_button_open") && (i = window.PayDiscount.left()); i = Math.round(i); j.ajax.post(c.util.buildAjaxRequest({
                            data: { realAmount: i }, url: c.constants.requestURL.getInstalmentInfoURL, success: function (h) {
                                h = JSON.parse(h); c.util.loadingFinish(); if (h.p.instalmentOptions != "unsupport" && h.p.instalmentOptions != "nodata") {
                                    e.settings.instalmentOptions = h.p.instalmentOptions; var i = e.settings.instalmentOptions[0], h = g(i), i = i.periods;
                                    b("#instalment-item-info").attr("periods", i); window.selectInstalmentPeriods = i; window.selectInstalmentContent = h; try { var j = h.split("|"), l = j[0] + " ", u = j[3].split(",")[0]; b("#instalment-item-info").html("<span>" + l + u + '</span><div class="jt_left"></div>').show() } catch (o) { b("#instalment-item-info").html("<span>" + h + '</span><div class="jt_left"></div>').show() } b(d).removeClass("discount_button_close"); b(d).addClass("discount_button_open")
                                } else if (h.p.instalmentOptions == "nodata") {
                                    e.settings.instalmentOptions =
                                    ""; b("#instalment-item-info").attr("periods", null); m.close(); c.util.toast(f.installmentOutOfAmount, "instalment_toast")
                                } else { e.settings.instalmentOptions = ""; b("#instalment-item-info").attr("periods", null); c.util.toast(f.notSupportInstalment, "instalment_toast"); m.close() }
                            }, error: function () { c.util.loadingFinish(); console && console.log("card validate\u5931\u8d25 \uff01") }
                        }))
                } else {
                    b(d).removeClass("discount_button_open"); b(d).addClass("discount_button_close"); window.selectInstalmentPeriods = ""; window.selectInstalmentContent =
                        ""; b("#instalment-item-info").empty(); b("#instalment-item-info").hide()
                }
            }); b("#instalment-item-info").unbind("click"); b("#instalment-item-info").click(function () {
                if (e.settings.supportInstalment && e.settings.instalmentOptions.length != 0) {
                    complete(); var d = ""; if (e.settings.supportInstalment == true) if (e.settings.instalmentOptions.length > 0) for (var i = 0; i < e.settings.instalmentOptions.length; i++) {
                        var h = e.settings.instalmentOptions[i]; g(h); var l = h.feeType, k = "", m = h.annualizedRateDisplay, k = h.periods == b("#instalment-item-info").attr("periods") ?
                            k + ("<li  periods=" + h.periods + '  content ="' + g(h) + '" ><div class="instalment-ul_left" ><a class="instalment-yh_select instalment-yh_select_yes"></a><span>' + h.periods + f.installment_term + '</span></div><div id="instalment-ul_right" class="instalment-ul_right"><p><span class="staging_wz">' + f.instalment_total_fee_desc + '</span><span class="staging_num">' + h.feeDisplay + f.rmb + '</span></p><p><span class="staging_wz">' + f.installment_averageAmount_fee_desc + '</span><span class="staging_num">' + h.averageDisplay +
                                f.rmb + "</span></p>" + (l == "allInfirst" ? '<p><span class="staging_wz">' + f.installment_firstAmount_fee_desc + '</span><span class="staging_num">' + h.firstDisplay + f.rmb + "</span></p>" : "") + '<p><span class="staging_wz">' + f.installment_payment_amount_desc + '</span><span class="staging_num">' + h.totleDisplay + f.rmb + "</span></p>" + (h.feeDisplay > 0 ? '<p><span class="staging_wz">' + f.installment_annualizedRate + '</span><span class="staging_num">' + m + "</span></p>" : "") + '</div><div style="clear:both"></div></li>') : k + ("<li  periods=" +
                                    h.periods + '  content ="' + g(h) + '" ><div class="instalment-ul_left" ><a class="instalment-yh_select "></a><span>' + h.periods + f.installment_term + '</span></div><div id="instalment-ul_right" class="instalment-ul_right"><p><span class="staging_wz">' + f.instalment_total_fee_desc + '</span><span class="staging_num">' + h.feeDisplay + f.rmb + '</span></p><p><span class="staging_wz">' + f.installment_averageAmount_fee_desc + '</span><span class="staging_num">' + h.averageDisplay + f.rmb + "</span></p>" + (l == "allInfirst" ? '<p><span class="staging_wz">' +
                                        f.installment_firstAmount_fee_desc + '</span><span class="staging_num">' + h.firstDisplay + f.rmb + "</span></p>" : "") + '<p><span class="staging_wz">' + f.installment_payment_amount_desc + '</span><span class="staging_num">' + h.totleDisplay + f.rmb + "</span></p>" + (h.feeDisplay > 0 ? '<p><span class="staging_wz">' + f.installment_annualizedRate + '</span><span class="staging_num">' + m + "</span></p>" : "") + '</div><div style="clear:both"></div></li>'), d = d + k
                    } else d = '<h1 id="topPoint-tips-msg">' + f.noInstalment + "</h1>"; b("#popInfo").empty().append(b(j.template.applyTpl(c.TPL.instalmentInfoHTML,
                        { staticURL: c.constants.staticURL, instalmentItemHTML: d, i18: f }))).show(); b(".echo_bank").show(); b("#instalment_list_ul li").unbind("click"); b("#instalment_list_ul li").click(function (c) {
                            c.preventDefault(); c = b(this).find(".instalment-yh_select"); if (c.hasClass("instalment-yh_select")) {
                                b(".instalment-yh_select").each(function () { b(this).removeClass("instalment-yh_select_yes") }); c.addClass("instalment-yh_select_yes"); window.selectInstalmentPeriods = parseInt(b(this).attr("periods")); window.selectInstalmentContent =
                                    b(this).attr("content"); setTimeout(function () { var c = window.selectInstalmentContent; try { var d = c.split("|"), e = d[0] + " ", f = d[3].split(",")[0]; b("#instalment-item-info").html("<span>" + e + f + '</span><div class="jt_left"></div>') } catch (g) { b("#instalment-item-info").html("<span>" + c + '</span><div class="jt_left"></div>') } b("#instalment-item-info").attr("periods", window.selectInstalmentPeriods); b("#instalment-item-info").show() }, 500)
                            }
                        }); b("#instalment_cancel_img,#confirm_title").unbind("click"); b("#instalment_cancel_img,#confirm_title").click(function () { b("#fq_close").hide() });
                    b(".fq_close").unbind("click"); b(".fq_close").click(function () {
                        if (window.selectInstalmentPeriods == void 0) { b("#fq_close").hide(); var c = document.getElementById("fqkg"), d = c.getElementsByTagName("span")[0], e = 1, f = function () { if (e < 20) { a = (21 - e) / 5; e = a + e; d.style.marginRight = e + "px"; b(c).removeClass("discount_button_open"); b(c).addClass("discount_button_close"); setTimeout(f, 10) } }; b("#instalment-item-info").attr("periods") || f() } else if (window.selectInstalmentPeriods == "") { b("#fq_close").hide(); b("#fqkg").trigger("click") } else {
                            var g =
                                window.selectInstalmentContent; b("#instalment-item-info").html("<span>" + g + '</span><div class="jt_left"></div>'); b("#instalment-item-info").attr("periods", window.selectInstalmentPeriods); b("#instalment-item-info").show(); b("#fq_close").hide()
                        }
                    })
                }
            })
        }
    }; m.close = function () { b("#fqkg").hasClass("discount_button_open") && b("#fqkg").trigger("click") }; m.instance = null; m.show = function (b) { m.options = b; m.instance = new m(b); return m.instance }; m.getPostData = function () {
        return c.constants.supportInstalment && b("#fqkg").hasClass("discount_button_open") ?
            window.selectInstalmentPeriods : null
    }; window.Instalment = m; j.util.getScript(window.mobileStaticUrl + getStaticMD5("/resources/upop_m/js/phone/phone.encryptpd.js"), function () { console && console.log("loading anatisy logging!!!!") }); return c
});
