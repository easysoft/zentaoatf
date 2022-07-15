import request from '@/utils/request';

const apiPath = 'settings';

export async function setLang(lang): Promise<any> {
    const params = {lang :lang}
    return request({
        url: `/${apiPath}/setLang`,
        method: 'GET',
        params
    });
}
