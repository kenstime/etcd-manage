<template>
    <div class="keys">
        <div class="path">
            <span @click="openKeyForPath({
                    key:'',
                    name:''
                })"><Icon type="ios-home-outline" /></span>
            <span v-show="showPathInput == false">
                <span v-for="(item,key) in paths" :key="key" @click="openKeyForPath(item)"><span v-if="key > 0 || paths.length == 1">/</span> {{ item.name }} </span>
            </span>
            
            <!-- 路径文本框 -->
            <Input v-show="showPathInput == true" @on-enter="enterPath" v-model="currentPath" placeholder="key路径，只能是路径不能是值" style="width: 300px;margin-top:-9px;font-size:28px;" />
            <span @click="showPathInput = !showPathInput" style="font-size:15px;margin-left:5px;">
                <span v-show="showPathInput == false">编辑</span>
                <span v-show="showPathInput == true">隐藏</span>
            </span>

            <Poptip style="float:right;margin-top:-3px;" placement="left-end"
                confirm
                title="确定删除选中的key？"
                @on-ok="delKeys">
            <Button  type="error">删除选中</Button>
            </Poptip>

            <Button type="success" @click="addKeyInfoModel = true;addKeyInfo = {kType:'KEY'}" style="float:right;margin-right:6px;">添加</Button>
            <!-- 展示方式切换 -->
            <RadioGroup v-model="listType" @on-change="selectionKeys=[]" type="button" style="float:right;margin-right:6px;">
                <Radio label="grid"><Icon type="md-grid" /></Radio>
                <Radio label="list"><Icon type="ios-list-box-outline" /></Radio>
            </RadioGroup>

        </div>

        <div v-if="keysList.length == 0">没有子级列表</div>

        <div class="lists">

            <!-- grid方式展示 -->
            <div class="key-list-main" v-if="listType == 'grid'">
                <div v-for="(item,key) in keysList" :key="key" class="key-list" @click="checkKey(item)">
                    <Checkbox class="checkbox" v-model="item.check"></Checkbox>
                    <Tooltip :content="item.key" placement="top-start">
                    <div @click.stop="openKey(item)">
                        <img v-if="item.is_dir==false" src="../assets/imgs/file.png" alt="file" class="key-icon">
                        <img v-if="item.is_dir==true" src="../assets/imgs/folder.png" alt="file" class="key-icon">
                    </div>
                    <div class="title">
                        {{item.path}}
                    </div>
                    </Tooltip>
                </div>
            </div>

            <!-- 表格形式展示 -->
            <div class="table-list-main" v-if="listType == 'list'">
                <Table border :columns="columnsKey" :data="keysList" @on-selection-change="selectionChangeKeys"></Table>
            </div>


            <div style="clear:both"></div>
            
        </div>

        <!-- 查看弹框 -->
        <Drawer
            :width="50"
            v-model="showKeyInfoModel"
            title="配置信息">
            <Form :model="showKeyInfo" :label-width="80">
                <FormItem label="Key" prop="key">
                    <Input v-model="showKeyInfo.key" disabled placeholder="配置地址"></Input>
                </FormItem>
                <FormItem label="Value" prop="value" >
                  <!-- <codemirror ref="code" v-model="showKeyInfo.value" :options="cmOption" style="line-height:20px;"></codemirror> -->
                    <!-- <Input v-model="showKeyInfo.value" type="textarea" :autosize="{minRows: 8,maxRows: 80}" placeholder="配置值"></Input> -->
                  <codemirror v-model="showKeyInfo.value" style="line-height:20px;"></codemirror>
                </FormItem>
                <FormItem>
                    <Button @click="saveKey" type="primary">保存</Button>
                    <Button @click="showKeyInfoModel = false" style="margin-left: 8px">关闭</Button>
                </FormItem>

            </Form>
        </Drawer>

        <!-- 添加弹框 -->
        <Drawer
            :width="50"
            v-model="addKeyInfoModel"
            title="添加KEY">
            <Form :model="addKeyInfo" :label-width="80">
                <FormItem label="Key" prop="key">
                    <Input v-model="addKeyInfo.key" placeholder="配置地址">
                        <span slot="prepend">{{currentPath}}/</span>
                    </Input>
                </FormItem>
                <FormItem label="kType" prop="kType">
                    <RadioGroup v-model="addKeyInfo.kType">
                        <Radio label="KEY"></Radio>
                        <Radio label="DIR"></Radio>
                    </RadioGroup>
                </FormItem>
                
                <FormItem label="Value" prop="value" v-show="addKeyInfo.kType == 'KEY'">
                  <!-- <codemirror v-model="addKeyInfo.value" :options="cmOption" style="line-height:20px;"></codemirror> -->
                    <!-- <Input v-model="addKeyInfo.value" type="textarea" :autosize="{minRows: 8,maxRows: 80}" placeholder="配置值"></Input> -->
                  <codemirror v-model="addKeyInfo.value" :options="{mode: 'javascript'}" style="line-height:20px;border: 1px solid #dcdee2"></codemirror>
                </FormItem>

                <FormItem>
                    <Button @click="addKey" type="primary">保存</Button>
                    <Button @click="addKeyInfoModel = false" style="margin-left: 8px">关闭</Button>
                </FormItem>

            </Form>
        </Drawer>

    </div>
</template>

<script>
require("codemirror/mode/javascript/javascript");
require("codemirror/mode/toml/toml");
require("codemirror/mode/yaml/yaml");
require("codemirror/mode/xml/xml");

export default {
  data() {
    return {
      showPathInput: false, // 是否显示路径文本框
      listType: "grid", // 显示方式
      paths: [
        {
          key: "/",
          name: ""
        }
      ], // 路径
      currentPath: "", // 当前key路径
      showKeyInfoModel: false, // 配置详情是否显示
      showKeyInfo: {}, // 要显示的配置值

      addKeyInfoModel: false, // 添加弹框
      addKeyInfo: {}, // 添加时对象

      selectionKeys: [], // 表格选中列表

      keysList: [],

      columnsKey: [
        // 表格形式展示表头
        {
          type: "selection",
          width: 60,
          align: "center"
        },
        {
          title: "TYPE",
          width: 70,
          render: (h, params) => {
            return h("img", {
              attrs: {
                src:
                  params.row.is_dir == true
                    ? "data:image/png;base64,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"
                    : "data:image/png;base64,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",
                style: "width:35px;height:35px"
              },
              on: {
                click: () => {
                  this.openKey(params.row);
                }
              }
            });
          }
        },
        {
          title: "NAME",
          key: "path"
        },
        {
          title: "KEY",
          key: "key"
        },
        {
          title: "操作",
          align: "center",
          render: (h, params) => {
            return h("div", [
              h(
                "Button",
                {
                  props: {
                    type: "primary",
                    size: "small"
                  },
                  style: {
                    marginRight: "5px"
                  },
                  on: {
                    click: () => {
                      this.openKey(params.row);
                    }
                  }
                },
                params.row.is_dir == true ? "打开" : "查看"
              ),
              h(
                "Poptip",
                {
                  props: {
                    confirm: true,
                    title: "确定删除"
                  },
                  on: {
                    "on-ok": () => {
                      this.delOneKey(params.row.key);
                    }
                  }
                },
                [
                  h(
                    "Button",
                    {
                      props: {
                        type: "error",
                        size: "small"
                      }
                    },
                    "删除"
                  )
                ]
              )
            ]);
          }
        }
      ],
      // 代码编辑器配置
      cmOption: {
        tabSize: 4,
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        mode: "text/javascript",
        lineWrapping: true,
        extraKeys: { Ctrl: "autocomplete" },
        theme: "monokai",
        autoRefresh: true
      }
    };
  },
  mounted() {
    this.getKeyList();
  },
  methods: {
    // 路径文本框回车
    enterPath() {
      this.openKeyForPath({
        key: this.currentPath,
        name: ""
      });
    },
    // 选中key
    checkKey(item) {
      item.check = !item.check;
    },
    // 打开key
    openKey(item) {
      console.log(item);
      if (item.is_dir == false) {
        this.showKeyInfo = item;
        this.showKeyInfoModel = true;
      } else {
        this.currentPath = item.key;
        let subarr = item.key.substring(1).split("/");
        this.paths.push({
          key: item.key,
          name: subarr[subarr.length - 1]
        });
        // console.log(this.paths);
        this.getKeyList();
      }
    },
    // 顶部路径打开目录
    openKeyForPath(item) {
      this.currentPath = item.key;
      let subKey = "";
      this.paths = [
        {
          key: "/",
          name: ""
        }
      ];
      if (item.key == "") {
        this.currentPath = "";
        console.log(this.paths);
        this.getKeyList();
        return;
      }

      item.key
        .substring(1)
        .split("/")
        .forEach(val => {
          subKey += "/" + val;
          this.paths.push({
            key: subKey,
            name: val
          });
        });
      console.log(this.paths);
      this.getKeyList();
    },

    // 添加key
    addKey() {
      if (this.addKeyInfo.key == "" || typeof this.addKeyInfo.key == "undefined") {
        this.$Message.warning("key不能为空");
        return;
      }
      this.addKeyInfo.key = this.currentPath + "/" + this.addKeyInfo.key;
      let uri = `/kv${this.addKeyInfo.key}`;
      console.log(this.addKeyInfo.kType);
      if (this.addKeyInfo.kType == "DIR") {
        uri += "?dir";
      }
      this.$http
        .post(uri, {
          value: this.addKeyInfo.value
        })
        .then(resp => {
          if (resp.data.err == "") {
            this.$Message.success("添加成功！");
            this.getKeyList();
            this.addKeyInfo = {};
            this.addKeyInfoModel = false;
          } else {
            this.$Message.error(resp.data.err);
          }
        });
    },

    saveKey() {
      this.$http
        .put(`/kv${this.showKeyInfo.key}`, {
          value: this.showKeyInfo.value
        })
        .then(resp => {
          if (resp.data.err == "") {
            this.$Message.success("保存成功！");
            this.showKeyInfoModel = false;
          } else {
            this.$Message.error(resp.data.err);
          }
        });
    },

    // 删除key
    delKeys() {
      this.keysList.forEach((val, index) => {
        // console.log(val._checked);
        if (val.check == true && this.listType == "grid") {
          this.delOneKey(val.key);
        }
      });
      if (this.listType == "list") {
        this.selectionKeys.forEach(val => {
          this.delOneKey(val.key);
        });
        this.selectionKeys = [];
      }
    },

    delOneKey(key) {
      this.$http.delete(`/kv${key}`).then(resp => {
        if (resp.data.err == "") {
          this.getKeyList();
        } else {
          this.$Message.error(resp.data.err);
        }
      });
    },

    // 表格选中行
    selectionChangeKeys(selections) {
      this.selectionKeys = selections;
      console.log(selections);
    },

    // 获取key列表
    getKeyList() {
      let k = this.currentPath == "" ? "/" : this.currentPath;
      this.$http.get(`/kv${k}?list`).then(resp => {
        if (resp.data.err == "") {
          this.keysList = resp.data.result;
          this.keysList.forEach((val,kkk) => {
            
            val.check = false;
            val._checked = false;
            let subarr = val.key.substring(1).split("/");
            val.path = subarr[subarr.length - 1];

            this.$set(this.keysList,kkk, val);
            console.log(val);
          });
        }else{
          this.$Message.error(resp.data.err);
        }
      });
    }
  }
};
</script>

<style scoped>
.keys {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
.keys .path {
  width: 100%;
  font-size: 24px;
  border-bottom: 1px solid #cecece;
}
.keys .path span {
  margin: 5px 0px 13px 0px;
  color: #333333;
  cursor: pointer;
}

.keys .lists {
  position: relative;
  width: 100%;
  height: 100vh;
  /* background-color: aqua; */
}

.keys .lists .key-list-main {
  position: absolute;
  left: 0;
  top: 0;
}

.keys .lists .key-list {
  position: relative;
  width: 110px;
  height: 110px;
  padding: 15px;
  margin: 10px;
  float: left;
  text-align: center;
  background-color: #eeeeee;
}

.keys .lists .key-list .title {
  font-size: 18px;
  width: 60px;
  overflow: hidden;
  /* text-overflow:ellipsis; */
  white-space: nowrap;
}

.keys .lists .key-list .checkbox {
  position: absolute;
  left: 10px;
  top: 10px;
}

.keys .lists .key-icon {
  width: 60px;
  height: 60px;
}

.keys .lists .table-list-main {
  margin-top: 10px;
}

</style>
