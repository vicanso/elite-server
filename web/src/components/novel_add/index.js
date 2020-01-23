import React from "react";
import {
  Spin,
  message,
  Card,
  Input,
  Select,
  Form,
  Button,
  Row,
  Col
} from "antd";

import "./novel_add.sass";
import { NOVEL_STATUS } from "../../vars";
import * as novelService from "../../services/novel";

const { TextArea } = Input;
const { Option } = Select;
const statusList = NOVEL_STATUS;

class NovelAdd extends React.Component {
  state = {
    submitting: false,
    status: 0,
    name: "",
    author: "",
    summary: "",
    titleRegexp: null,
    chapters: null
  };
  async handleSubmit(e) {
    e.preventDefault();
    const { name, author, summary, chapters, submitting, status } = this.state;
    if (submitting) {
      return;
    }
    if (!name || !author || !chapters) {
      message.warning("名称、作者以及章节不能为空");
      return;
    }
    this.setState({
      submitting: true
    });
    try {
      await novelService.add({
        status,
        name,
        author,
        summary,
        chapters
      });
      message.info("已成功添加该书籍.");
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        submitting: false
      });
    }
  }
  formatChapters(data) {
    const { titleRegexp } = this.state;
    if (!titleRegexp) {
      message.warning("请先输入标题的正则表达式");
      return;
    }
    const arr = data
      .trim()
      .split("\n")
      .map(item => item.trim());
    const result = [];
    let content = [];
    let title = "";
    arr.forEach(item => {
      if (item.match(titleRegexp)) {
        if (title) {
          result.push({
            title,
            content: content.join("\n")
          });
        }
        title = item;
        content = [];
      } else {
        if (item) {
          content.push(item);
        }
      }
    });
    this.setState({
      chapters: result
    });
  }
  renderChapters() {
    const { chapters } = this.state;
    if (!chapters) {
      return;
    }
    const items = chapters.map(item => {
      return (
        <li>
          <h5>{item.title}</h5>
          <p>{item.content.substring(0, 30) + "..."}</p>
        </li>
      );
    });
    return <ul className="chaptersPreview">{items}</ul>;
  }
  renderStatusList() {
    return statusList.map(item => (
      <Option key={item} value={item}>
        {item}
      </Option>
    ));
  }
  render() {
    const colSpan = 8;
    const { submitting } = this.state;
    return (
      <div className="NovelAdd">
        <Card title="上传书籍" size="small">
          <Spin spinning={submitting}>
            <Form onSubmit={this.handleSubmit.bind(this)}>
              <Row gutter={24}>
                <Col span={colSpan}>
                  <Form.Item label="名称">
                    <Input
                      type="text"
                      onChange={e => {
                        this.setState({
                          name: e.target.value
                        });
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={colSpan}>
                  <Form.Item label="作者">
                    <Input
                      type="text"
                      onChange={e => {
                        this.setState({
                          author: e.target.value
                        });
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={colSpan}>
                  <Form.Item label="状态">
                    <Select
                      placeholder="请选择书籍状态"
                      onChange={value => {
                        this.setState({
                          status: statusList.indexOf(value)
                        });
                      }}
                    >
                      {this.renderStatusList()}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={24}>
                  <Form.Item label="摘要">
                    <TextArea
                      rows={6}
                      onChange={e => {
                        this.setState({
                          summary: e.target.value
                        });
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item label="章节内容">
                    <Input
                      type="text"
                      placeholder="请输入标题正则表达式：第\d+章"
                      onChange={e => {
                        try {
                          this.setState({
                            titleRegexp: new RegExp(e.target.value)
                          });
                        } catch (err) {}
                      }}
                    />
                    <TextArea
                      rows={12}
                      onChange={e => {
                        this.formatChapters(e.target.value);
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item label="章节预览">
                    {this.renderChapters()}
                  </Form.Item>
                </Col>
                <Col span={24}>
                  <Button className="submit" type="primary" htmlType="submit">
                    确认
                  </Button>
                </Col>
              </Row>
            </Form>
          </Spin>
        </Card>
      </div>
    );
  }
}

export default NovelAdd;
