/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { IAM_VIEWS } from '@/dictionary/iam-auth'
import CombineRequest from '@/api/combine-request.js'
import { foreignkey } from '@/filters/formatter.js'
import instanceService from '@/service/instance/instance'
import hostSearchService from '@/service/host/search'
import projectService from '@/service/project'
import businessSetService from '@/service/business-set/index.js'
import fieldTemplateService from '@/service/field-template'
import {
  BUILTIN_MODELS,
  BUILTIN_MODEL_PROPERTY_KEYS
} from '@/dictionary/model-constants.js'

const requestConfigBase = key => ({
  requestId: `permission_${key}`,
  fromCache: true
})

async function getBusinessList(vm) {
  // 使用`biz/search/${rootGetters.supplierAccount}`需要鉴权，从而使用biz/simplify
  const url = 'biz/simplify'
  const data = await vm.$http.get(`${url}?sort=bk_biz_id`, { ...requestConfigBase(url) })
  return data.info || []
}

async function getResourceDirectoryList(vm) {
  const action = 'resourceDirectory/getDirectoryList'
  let directoryList = vm.$store.getters['resourceHost/directoryList']
  if (!directoryList.length) {
    const res = await vm.$store.dispatch(action, { params: {}, config: { ...requestConfigBase(action) } })
    directoryList = res.info || []
  }
  return directoryList
}

export const IAM_VIEWS_INST_NAME = {
  [IAM_VIEWS.MODEL_GROUP](vm, id) {
    const classifications = vm.$store.getters['objectModelClassify/classifications']
    const value = (classifications.find(item => item.id === Number(id)) || {}).bk_classification_name
    return Promise.resolve(value)
  },
  [IAM_VIEWS.MODEL](vm, id) {
    const models = vm.$store.getters['objectModelClassify/models']
    const value = (models.find(item => item.id === Number(id)) || {}).bk_obj_name
    return Promise.resolve(value)
  },
  async [IAM_VIEWS.INSTANCE](vm, id, relations) {
    const models = vm.$store.getters['objectModelClassify/models']
    const objId = (models.find(item => item.id === Number(relations[0][1])) || {}).bk_obj_id
    const inst = await instanceService.findOne({
      bk_obj_id: objId,
      bk_inst_id: Number(id),
      config: { ...requestConfigBase(`find_instance_${id}`) }
    })
    return inst ? inst.bk_inst_name : id
  },
  async [IAM_VIEWS.PROJECT](vm, id) {
    const project = await projectService.findOne({
      id: Number(id)
    }, { ...requestConfigBase(`find_project_${id}`) })
    return project ? project.bk_project_name : id
  },
  [IAM_VIEWS.INSTANCE_MODEL](vm, id) {
    const models = vm.$store.getters['objectModelClassify/models']
    const value = (models.find(item => item.id === Number(id)) || {}).bk_obj_name
    return Promise.resolve(value)
  },
  async [IAM_VIEWS.CUSTOM_QUERY](vm, id, relations) {
    const bizId = Number(relations[0][1])
    const action = 'dynamicGroup/details'
    const details = await vm.$store.dispatch(action, {
      bizId,
      id,
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    const value = details.name
    return value
  },
  async [IAM_VIEWS.BIZ](vm, id) {
    const list = await getBusinessList(vm)
    const business = list.find(business => business.bk_biz_id === Number(id))
    return business.bk_biz_name
  },
  async [IAM_VIEWS.BIZ_SET](vm, id) {
    const { info: list = [] } = await businessSetService.getAll(requestConfigBase(`biz_set_${id}`))
    const MODEL_ID_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID
    const MODEL_NAME_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME
    const businessSet = list.find(item => item[MODEL_ID_KEY] === Number(id))
    return businessSet[MODEL_NAME_KEY]
  },
  async [IAM_VIEWS.BIZ_FOR_HOST_TRANS](vm, id) {
    const list = await getBusinessList(vm)
    const business = list.find(business => business.bk_biz_id === Number(id))
    return business.bk_biz_name
  },
  async [IAM_VIEWS.HOST](vm, id) {
    const action = 'hostSearch/searchHost'

    const result = await CombineRequest.setup(action, async (data) => {
      const hostIdList = data.map(Number)
      const hostCondition = {
        field: 'bk_host_id',
        operator: '$in',
        value: hostIdList
      }
      const params = {
        bk_biz_id: -1,
        condition: ['biz', 'set', 'module', 'host'].map(model => ({
          bk_obj_id: model,
          condition: model === 'host' ? [hostCondition] : [],
          fields: []
        })),
        ip: { flag: 'bk_host_innerip', exact: 1, data: [] }
      }

      // 这个的前置场景是主机的编辑，无论编辑的业务主机/资源池主机，在进入到操作入口时应该具体了入口权限，因此这里使用无需鉴权的查询
      const { info } = await hostSearchService.getHosts({
        params,
        config: { ...requestConfigBase(`${action}${hostIdList.join('')}`) }
      })

      return info
    }).add(id)

    const { host } = result.find(({ host }) => host.bk_host_id === Number(id)) || {}
    return `${foreignkey(host.bk_cloud_id)}: ${host.bk_host_innerip}`
  },
  async [IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY](vm, id) {
    const directoryList = await getResourceDirectoryList(vm)
    const directory = directoryList.find(directory => directory.bk_module_id === Number(id)) || {}
    return directory.bk_module_name
  },
  async [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY](vm, id) {
    const directoryList = await getResourceDirectoryList(vm)
    const directory = directoryList.find(directory => directory.bk_module_id === Number(id)) || {}
    return directory.bk_module_name
  },
  async [IAM_VIEWS.ASSOCIATION_TYPE](vm, id) {
    const action = 'objectAssociation/searchAssociationType'
    const { info: associationList } = await vm.$store.dispatch(action, {
      params: {},
      config: { ...requestConfigBase(action) }
    })
    const asst = associationList.find(asst => asst.id === Number(id))
    return asst.bk_asst_name
  },
  async [IAM_VIEWS.EVENT_PUSHING](vm, id) {
    const action = 'eventSub/searchSubscription'
    const res = await vm.$store.dispatch(action, {
      bkBizId: 0,
      params: {
        page: {
          start: 0,
          limit: 1
        },
        condition: {
          subscription_id: Number(id)
        }
      },
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    const subscription = res.info[0] || {}
    return subscription.subscription_name
  },
  async [IAM_VIEWS.SERVICE_TEMPLATE](vm, id) {
    const action = 'serviceTemplate/findServiceTemplate'
    const serviceTemplate = await vm.$store.dispatch(action, {
      id: Number(id),
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    const template = serviceTemplate.template || {}
    return template.name
  },
  async [IAM_VIEWS.SET_TEMPLATE](vm, id) {
    const action = 'setTemplate/getSetTemplates'
    const res = await vm.$store.dispatch(action, {
      bizId: vm.$store.getters['objectBiz/bizId'],
      params: {
        set_template_ids: [id].map(Number)
      },
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    const data = res.info[0] || {}
    const setTemplate = data.set_template || {}
    return setTemplate.name
  },
  async [IAM_VIEWS.CLOUD_AREA](vm, id) {
    const action = 'cloud/area/findMany'
    const res = await vm.$store.dispatch(action, {
      params: {
        condition: {
          bk_cloud_id: Number(id)
        }
      },
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    const data = res.info[0] || {}
    return data.bk_cloud_name
  },
  async [IAM_VIEWS.CLOUD_ACCOUNT](vm, id) {
    const action = 'cloud/account/findOne'
    const account = await vm.$store.dispatch(action, {
      id: Number(id),
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    return account.bk_account_name
  },
  async [IAM_VIEWS.CLOUD_RESOURCE_TASK](vm, id) {
    const action = 'cloud/resource/findOneTask'
    const task = await vm.$store.dispatch(action, {
      id: Number(id),
      config: { ...requestConfigBase(`${action}${id}`) }
    })
    return task?.bk_task_name
  },
  async [IAM_VIEWS.FIELD_TEMPLATE](vm, id) {
    const res = await fieldTemplateService.findById(id, requestConfigBase('field_template'))
    return res.name
  }
}
