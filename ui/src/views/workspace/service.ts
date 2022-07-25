import request from '@/utils/request';
import { QueryParams } from '@/types/data.d';
import { checkProxy } from "@/views/proxy/service";

const apiPath = 'workspaces';

export async function query(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function listByProduct(productId: number): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listByProduct`,
        method: 'get',
        params,
    });
}

export async function get(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`
    });
}

export async function save(params: any): Promise<any> {
    return request({
        url: `/${apiPath}` + (params.id ? `/${params.id}` : ''),
        method: params.id? 'PUT': 'POST',
        data: params,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

export async function uploadToProxy(params: any): Promise<any> {
    return request({
        url: `/${apiPath}/uploadToProxy`,
        method: 'POST',
        data: params,
    });
}

export async function autoSelectProxy(workspace) {
  const proxies = workspace.proxies.split(',');
  const handleList = [] as any;
  let localIndex = proxies.length;
  proxies.forEach((proxy, index) => {
    if (proxy > 0) {
      handleList.push(
        checkProxy(proxy)
      );
    }else{
        localIndex = index;
    }
  })
  const resp = await Promise.all(handleList);
  let proxyPath = '';
  resp.forEach((item:any, index) => {
    if(proxyPath == '' && item.data.status == 'ok' && index < localIndex){
        proxyPath = item.data.path;
    }
  })
  return proxyPath ? proxyPath : 'local';
}