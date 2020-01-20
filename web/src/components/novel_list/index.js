import React from "react";
import {
  message,
  Table,
  Card,
  Spin,
  Input,
  Select,
  Form,
  Button,
  Row,
  Col
} from "antd";
import moment from "moment";

import "./novel_list.sass";
import * as novelService from "../../services/novel";
import { TIME_FORMAT } from "../../vars";

const { Search, TextArea } = Input;
const { Option } = Select;
const allStatus = "所有";
const statusList = ["未知", "未完结", "完结", "下架"];
const editMode = "edit";

class NovelList extends React.Component {
  state = {
    mode: "",
    current: null,
    updateData: null,
    loading: false,
    submitting: false,
    keyword: "",
    status: "",
    order: "",
    pagination: {
      current: 1,
      pageSize: 10,
      total: 0
    },
    novels: null
  };
  reset() {
    const pagination = Object.assign(
      { ...this.state.pagination },
      {
        current: 1,
        total: 0
      }
    );
    this.setState({
      novels: null,
      pagination
    });
  }
  async list() {
    const { loading, pagination, keyword, status, order } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      const offset = (pagination.current - 1) * pagination.pageSize;
      const data = await novelService.list({
        order,
        status,
        keyword,
        limit: pagination.pageSize,
        offset
      });
      const updateData = {
        novels: data.novels
      };
      if (data.count >= 0) {
        updateData.pagination = Object.assign(
          { ...pagination },
          {
            total: data.count
          }
        );
      }
      this.setState(updateData);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  async handleSubmit() {
    const { current, updateData, submitting } = this.state;
    if (submitting) {
      return;
    }
    this.setState({
      submitting: true
    });
    try {
      await novelService.updateByID(current.id, updateData);
      Object.assign(current, updateData);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        submitting: false,
        mode: ""
      });
    }
  }
  async updateCover(imageURL) {
    const { current } = this.state;
    try {
      await novelService.updateCoverByID(current.id, imageURL);
      message.info("更新封面成功！");
    } catch (err) {
      message.error(err.message);
    }
  }
  componentDidMount() {
    this.list();
  }
  renderStatusList(includeAll) {
    const arr = statusList.slice(1);
    if (includeAll) {
      arr.unshift(allStatus);
    }
    return arr.map(item => (
      <Option key={item} value={item}>
        {item}
      </Option>
    ));
  }
  renderNovelList() {
    const { loading, mode } = this.state;
    if (mode === editMode) {
      return;
    }
    return (
      <div>
        <Card title="书籍搜索" size="small">
          <Spin spinning={loading}>
            <div className="filter">
              <Select
                defaultValue={allStatus}
                className="statusList"
                placeholder="请选择书籍状态"
                onChange={value => {
                  const index = statusList.indexOf(value);
                  let status = "";
                  if (index >= 0) {
                    status = `${index}`;
                  }
                  this.setState(
                    {
                      status
                    },
                    () => {
                      this.reset();
                      this.list();
                    }
                  );
                }}
              >
                {this.renderStatusList(true)}
              </Select>
              <Search
                className="keyword"
                placeholder="请输入关键字"
                onSearch={keyword => {
                  this.setState(
                    {
                      keyword
                    },
                    () => {
                      this.reset();
                      this.list();
                    }
                  );
                }}
                enterButton
              />
            </div>
          </Spin>
        </Card>
        {this.renderTable()}
      </div>
    );
  }
  renderTable() {
    const { novels, pagination } = this.state;
    const columns = [
      {
        title: "名称",
        dataIndex: "name",
        key: "name"
      },
      {
        title: "作者",
        dataIndex: "author",
        key: "author"
      },
      {
        title: "分级",
        dataIndex: "grading",
        key: "grading",
        sorter: true
      },
      {
        title: "状态",
        dataIndex: "status",
        key: "status",
        render: status => {
          return statusList[status || 0];
        }
      },
      {
        title: "封面",
        render: (text, record) => {
          return (
            <img
              alt="cover"
              src={`/novels/v1/${record.id}/cover?output=webp&width=40&nocache=true`}
            />
          );
        }
      },
      {
        title: "更新时间",
        dataIndex: "updatedAt",
        key: "updatedAt",
        sorter: true,
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      },
      {
        title: "操作",
        key: "op",
        width: "100px",
        render: (text, record) => {
          return (
            <a
              href="/update"
              onClick={e => {
                e.preventDefault();
                this.setState({
                  updateData: {},
                  current: record,
                  mode: editMode
                });
              }}
            >
              更新
            </a>
          );
        }
      }
    ];

    return (
      <Table
        rowKey={"id"}
        className="novels"
        dataSource={novels}
        columns={columns}
        pagination={pagination}
        onChange={(pagination, filters, sorter) => {
          const { field, order } = sorter;
          let orderBy = "";
          if (field && order) {
            if (order === "descend") {
              orderBy = `-${field}`;
            } else {
              orderBy = field;
            }
          }
          this.setState(
            {
              pagination: { ...pagination },
              order: orderBy
            },
            () => {
              this.list();
            }
          );
        }}
      />
    );
  }
  renderEditor() {
    const { mode, current, updateData } = this.state;
    if (mode !== editMode) {
      return;
    }
    const colSpan = 8;
    return (
      <Card title="更新书籍信息" size="small">
        <Form onSubmit={this.handleSubmit.bind(this)}>
          <Row gutter={24}>
            <Col span={colSpan}>
              <Form.Item label="名称">
                <Input disabled defaultValue={current.name} />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="状态">
                <Select
                  defaultValue={statusList[current.status]}
                  placeholder="请选择书籍状态"
                  onChange={value => {
                    updateData.status = statusList.indexOf(value);
                  }}
                >
                  {this.renderStatusList(false)}
                </Select>
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="分级">
                <Input
                  type="number"
                  defaultValue={current.grading}
                  onChange={e => {
                    updateData.grading = e.target.valueAsNumber;
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="封面">
                <Input
                  type="text"
                  placeholder="请输入封面图片的抓取地址"
                  onChange={e => {
                    this.updateCover(e.target.value);
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item label="摘要">
                <TextArea
                  defaultValue={current.summary}
                  rows={6}
                  onChange={e => {
                    updateData.summary = e.target.value;
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Button className="submit" type="primary" htmlType="submit">
                更新
              </Button>
            </Col>
            <Col span={colSpan}>
              <Button
                className="back"
                onClick={() => {
                  this.setState({
                    mode: ""
                  });
                }}
              >
                返回
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    );
  }
  render() {
    return (
      <div className="NovelList">
        {this.renderNovelList()}
        {this.renderEditor()}
      </div>
    );
  }
}

export default NovelList;
