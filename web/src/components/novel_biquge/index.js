import React from "react";
import {
  message,
  Card,
  Input,
  Button,
  Row,
  Col
} from "antd";

import "./novel_biquge.sass";

import * as novelService from "../../services/novel";

class NovelBiQuGe extends React.Component {
  state = {
    max: 0
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
  render() {
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
      </div>
    );
  }
}

export default NovelBiQuGe;
