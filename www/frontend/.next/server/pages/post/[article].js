"use strict";
(() => {
var exports = {};
exports.id = 803;
exports.ids = [803];
exports.modules = {

/***/ 124:
/***/ ((module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.a(module, async (__webpack_handle_async_dependencies__) => {
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "getStaticPaths": () => (/* binding */ getStaticPaths),
/* harmony export */   "getStaticProps": () => (/* binding */ getStaticProps),
/* harmony export */   "default": () => (__WEBPACK_DEFAULT_EXPORT__)
/* harmony export */ });
/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(997);
/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__);
/* harmony import */ var js_yaml__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(626);
/* harmony import */ var gray_matter__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(76);
/* harmony import */ var gray_matter__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(gray_matter__WEBPACK_IMPORTED_MODULE_2__);
/* harmony import */ var marked__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(974);
var __webpack_async_dependencies__ = __webpack_handle_async_dependencies__([marked__WEBPACK_IMPORTED_MODULE_3__, js_yaml__WEBPACK_IMPORTED_MODULE_1__]);
([marked__WEBPACK_IMPORTED_MODULE_3__, js_yaml__WEBPACK_IMPORTED_MODULE_1__] = __webpack_async_dependencies__.then ? await __webpack_async_dependencies__ : __webpack_async_dependencies__);




const getStaticPaths = async ()=>{
    const res1 = await fetch("https://raw.githubusercontent.com/gotti/blog/main/contents/blog.yaml").then((res)=>res.blob()
    ).then((blob)=>blob.text()
    );
    console.log(res1);
    const y = js_yaml__WEBPACK_IMPORTED_MODULE_1__.load(res1)["posts"];
    const paths = y.map((post)=>{
        console.log(post);
        return post.slice(1);
    });
    return {
        paths,
        fallback: false
    };
};
const getStaticProps = async ({ params  })=>{
    console.log("params", params);
    const a = await fetch("https://raw.githubusercontent.com/gotti/blog/main/contents/post/" + params.article + "/index.md").then((res)=>res.blob()
    ).then((blob)=>blob.text()
    );
    console.log("https://raw.githubusercontent.com/gotti/blog/main/contents/post/" + params.article + "/index.md");
    const b = gray_matter__WEBPACK_IMPORTED_MODULE_2___default()(a);
    return {
        props: {
            page: b.content
        }
    };
};
const Article = ({ page  })=>{
    return(/*#__PURE__*/ react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.Fragment, {
        children: /*#__PURE__*/ react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx("div", {
            className: "postBody",
            children: /*#__PURE__*/ react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx("div", {
                dangerouslySetInnerHTML: {
                    __html: (0,marked__WEBPACK_IMPORTED_MODULE_3__.marked)(page)
                }
            })
        })
    }));
};
/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (Article);

});

/***/ }),

/***/ 76:
/***/ ((module) => {

module.exports = require("gray-matter");

/***/ }),

/***/ 997:
/***/ ((module) => {

module.exports = require("react/jsx-runtime");

/***/ }),

/***/ 626:
/***/ ((module) => {

module.exports = import("js-yaml");;

/***/ }),

/***/ 974:
/***/ ((module) => {

module.exports = import("marked");;

/***/ })

};
;

// load runtime
var __webpack_require__ = require("../../webpack-runtime.js");
__webpack_require__.C(exports);
var __webpack_exec__ = (moduleId) => (__webpack_require__(__webpack_require__.s = moduleId))
var __webpack_exports__ = (__webpack_exec__(124));
module.exports = __webpack_exports__;

})();