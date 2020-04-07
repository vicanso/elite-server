import React from "react";
import { message, Card, Input, Button, Row, Col, Table, Spin } from "antd";
import moment from "moment";

import "./novel_biquge.sass";

import * as novelService from "../../services/novel";
import { TIME_FORMAT } from "../../vars";

class NovelBiQuGe extends React.Component {
  state = {
    max: 0,
    keyword: "",
    loading: false,
    pagination: {
      current: 1,
      pageSize: 10,
      total: 0
    },
    novels: null
  };
  async syncBiQuGe() {
    const { max } = this.state;
    if (!max || max < 0) {
      message.warning("同步最大ID不能为空或小于0");
      return;
    }
    try {
      novelService.syncBiQuGe(max);
      message.info("同步任务已成功执行");
    } catch (err) {
      message.error(err.message);
    }
  }
  async list() {
    const { loading, pagination, keyword } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      const offset = (pagination.current - 1) * pagination.pageSize;
      const data = await novelService.listBiQuGe({
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
  async syncNovel(id) {
    const { loading } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      await novelService.syncNovel({
        source: "biquge",
        bookID: id
      });
      message.info("同步成功！");
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  componentDidMount() {
    this.list();
  }
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
        title: "书籍ID",
        dataIndex: "bookID",
        key: "bookID"
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
              href="/sync"
              onClick={e => {
                e.preventDefault();
                this.syncNovel(record.bookID);
              }}
            >
              同步
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
          this.setState(
            {
              pagination: { ...pagination }
            },
            () => {
              this.list();
            }
          );
        }}
      />
    );
  }
  render() {
    const { loading } = this.state;
    return (
      <div className="NovelBiQuGe">
        <Card title="同步书籍" size="small">
          <Row gutter={24}>
            <Col span={12}>
              <Input
                type="number"
                placeholder="请输入同步的最大ID"
                onChange={e => {
                  this.setState({
                    max: e.target.valueAsNumber
                  });
                }}
              />
            </Col>
            <Col span={12}>
              <Button
                onClick={e => {
                  this.syncBiQuGe();
                }}
              >
                确认同步
              </Button>
            </Col>
          </Row>
        </Card>
        <Card title="书籍列表" size="small" className="novels">
          <Row gutter={24}>
            <Col span={12}>
              <Input
                type="text"
                placeholder="请输入关键字"
                allowClear
                onChange={e => {
                  this.setState({
                    keyword: e.target.value
                  });
                }}
              />
            </Col>
            <Col span={12}>
              <Button
                onClick={e => {
                  this.reset();
                  this.list();
                }}
              >
                查询
              </Button>
            </Col>
          </Row>
          <Spin spinning={loading}>{this.renderTable()}</Spin>
        </Card>
      </div>
    );
  }
}

export default NovelBiQuGe;
