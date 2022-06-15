
export function resizeWidth(mainId: string, leftId: string, resizeId: string, contentId: string,
                            leftMin: number, rightMin: number): boolean {
    const main = document.getElementById(mainId) as any;
    const left = document.getElementById(leftId) as any;
    const resize = document.getElementById(resizeId) as any;
    const content = document.getElementById(contentId) as any;

    if (!resize) return false

    // 鼠标按下事件
    resize.onmousedown = function (e) {
        //色彩高亮
        resize.classList.add('active');
        const startX = e.clientX;

        // 鼠标拖动事件
        document.onmousemove = function (e) {
            resize.left = startX;
            const endX = e.clientX;
            let moveLen = resize.left + (endX - startX); // （endx-startx）=挪动的间隔。resize.left+挪动的间隔=右边区域最初的宽度
            const maxT = main.clientWidth - resize.offsetWidth; // 容器宽度 - 右边区域的宽度 = 左边区域的宽度
            if (moveLen < leftMin) moveLen = leftMin; // 右边区域的最小宽度
            if (moveLen > maxT - rightMin) moveLen = maxT - rightMin; //左边区域最小宽度
            resize.style.left = moveLen; // 设置左侧区域的宽度

            left.style.width = (moveLen / document.body.clientWidth) * 100 + '%';
            content.style.width = ((main.clientWidth - moveLen) / document.body.clientWidth - 0.008) * 100 + '%';
        };

        // 鼠标松开事件
        document.onmouseup = function (evt) {
            resize.classList.remove('active'); //色彩复原

            document.onmousemove = null;
            document.onmouseup = null;
            resize.releaseCapture && resize.releaseCapture(); //当你不在须要持续取得鼠标音讯就要应该调用ReleaseCapture()开释掉
        };

        resize.setCapture && resize.setCapture(); //该函数在属于以后线程的指定窗口里设置鼠标捕捉
        return false;
    };

    return true
}

export function resizeHeight(contentId: string, topId: string, splitterId: string, bottomId: string,
                             topMin: number, bottomMin: number, gap: number): boolean {
    const content = document.getElementById(contentId) as any;
    const top = document.getElementById(topId) as any;
    const splitter = document.getElementById(splitterId) as any;
    const bottom = document.getElementById(bottomId) as any;

    if (!splitter) {
        return false
    }

    // 鼠标按下事件
    splitter.onmousedown = function (e) {
        //色彩高亮
        splitter.classList.add('active');
        const startY = e.clientY - gap;

        // 鼠标拖动事件
        document.onmousemove = function (e) {
            splitter.top = startY;
            const endY = e.clientY - gap;
            let moveLen = splitter.top + (endY - startY); // （endY-startY）=挪动的间隔。splitter.top+挪动的间隔=上边区域最初的高度
            const maxT = content.clientHeight - splitter.offsetHeight; // 容器高度 - 下边区域的宽度 = 上边区域的宽度
            if (moveLen < topMin) moveLen = topMin; // 下边区域的最小宽度
            if (moveLen > maxT - bottomMin) moveLen = maxT - bottomMin; //上边区域最小高度
            splitter.style.top = moveLen; // 设置上边区域的高度

            top.style.height = (moveLen / content.clientHeight) * 100 + '%';
            bottom.style.height = ((content.clientHeight - moveLen) / content.clientHeight - 0.008) * 100 + '%';
        };

        // 鼠标松开事件
        document.onmouseup = function (evt) {
            splitter.classList.remove('active'); //色彩复原

            document.onmousemove = null;
            document.onmouseup = null;
            splitter.releaseCapture && splitter.releaseCapture(); //当你不在须要持续取得鼠标音讯就要应该调用ReleaseCapture()开释掉
        };

        splitter.setCapture && splitter.setCapture(); //该函数在属于以后线程的指定窗口里设置鼠标捕捉
        return false;
    };

    return true
}

export function PrefixZero(num: number, n: number): string {
    return (Array(n).join('0') + num).slice(-n);
}
export function PrefixSpace(num: number, n: number): string {
    return (Array(n).join(' ') + num).slice(-n);
}

export function scroll(id: string): void {
    const elem = document.getElementById(id)
    if (elem) {
        setTimeout(function(){
            elem.scrollTop = elem.scrollHeight + 100;
        },300);
    }
}

export function hasClass( elements, cName ){
    if (!elements) return false
    return !!elements.className.match( new RegExp( "(\\s|^)" + cName + "(\\s|$)") )
}
export function addClass( elements, cName ){
    if (!elements) return
    if( !hasClass( elements,cName ) ){
        elements.className += " " + cName
    }
}
export function removeClass( elements, cName ){
    if (!elements) return
    if( hasClass( elements,cName ) ){
        elements.className = elements.className.replace( new RegExp( "(\\s|^)" + cName + "(\\s|$)" ), " " )
    }
}

export function jsonStrDef(obj) {
    const msg = JSON.stringify(obj)
    return msg
}

export function getContextMenuStyle(x, y, height) {
    let top = y
    if (y + height > document.body.clientHeight)
        top = document.body.clientHeight - height
    const menuStyle = {
        zIndex: 9,
        position: 'fixed',
        left: `${x + 10}px`,
        top: `${top}px`,
    }

    return menuStyle
}
