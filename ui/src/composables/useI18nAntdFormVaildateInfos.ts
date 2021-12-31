/**
 * 重置 Antd Form VaildateInfos 为 I18n
 * @author LiQingSong
 */
import { computed, ComputedRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { validateInfos } from 'ant-design-vue/lib/form/useForm';

export default function useI18nAntdFormVaildateInfos(infos: validateInfos): ComputedRef<validateInfos> {
    const{ t } = useI18n();

    const infosNew = computed<validateInfos>(() => {
        const vInfos: validateInfos  = {};  
        for (const index in infos) {
            vInfos[index] = JSON.parse(JSON.stringify(infos[index]));           
            if(vInfos[index] && vInfos[index]['help']) {
                vInfos[index]['help'] = vInfos[index]['help'].map((item: any)=> typeof(item)=='string' ? t(item) : item.map((item2: any)=> item2 ? t(item2):''));
            }
        }
        return vInfos;
    });

    return infosNew;
}