
export function resizeWidth(mainId, leftId, resizeId, contentId) {
    const main = document.getElementById(mainId) as any;
    const left = document.getElementById(leftId) as any;
    const resize = document.getElementById(resizeId) as any;
    const content = document.getElementById(contentId) as any;

    // 鼠标按下事件
    resize.onmousedown = function (e) {
        //色彩扭转揭示
        resize.classList.add('active');
        const startX = e.clientX;
        // 鼠标拖动事件
        document.onmousemove = function (e) {
            resize.left = startX;
            const endX = e.clientX;
            let moveLen = resize.left + (endX - startX); // （endx-startx）=挪动的间隔。resize.left+挪动的间隔=右边区域最初的宽度
            const maxT = main.clientWidth - resize.offsetWidth; // 容器宽度 - 右边区域的宽度 = 左边区域的宽度
            if (moveLen < 32) moveLen = 600; // 右边区域的最小宽度为32px
            if (moveLen > maxT - 150) moveLen = maxT - 150; //左边区域最小宽度为150px
            resize.style.left = moveLen; // 设置左侧区域的宽度

            left.style.width = (moveLen / document.body.clientWidth) * 100 + '%';
            content.style.width = ((main.clientWidth - moveLen) / document.body.clientWidth - 0.008) * 100 + '%';
        };
        // 鼠标松开事件
        document.onmouseup = function (evt) {
            //色彩复原
            resize.classList.remove('active');
            document.onmousemove = null;
            document.onmouseup = null;
            resize.releaseCapture && resize.releaseCapture(); //当你不在须要持续取得鼠标音讯就要应该调用ReleaseCapture()开释掉
        };

        resize.setCapture && resize.setCapture(); //该函数在属于以后线程的指定窗口里设置鼠标捕捉
        return false;
    };
}