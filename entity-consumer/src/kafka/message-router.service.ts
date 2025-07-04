/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { Injectable } from '@nestjs/common';
import { MessageType } from '@protobuf-ts/runtime';
import { Entity } from 'gen/entity';
// import * as protobuf from 'protobufjs';

type ProtoTypes = Entity;

@Injectable()
export class MessageRouterService {
  private schemaMap: Record<string, MessageType<ProtoTypes>> = {};
  private handlerMap: Record<
    string,
    (data: ProtoTypes | null) => Promise<void>
  > = {};

  constructor() {} // TODO: add services here

  loadSchemas() {
    this.schemaMap['entity-topic'] = Entity;

    // this.handlerMap['entity-topic'] = this.userWorkflow.handleUserCreated.bind(
    //   this.userWorkflow,
    // );
  }

  async handle(topic: string, raw: Buffer) {
    const schema: MessageType<Entity> = this.schemaMap[topic];
    const handler = this.handlerMap[topic];

    if (!schema || !handler) {
      throw new Error(`No handler/schema for topic: ${topic}`);
    }

    let decoded: ProtoTypes | null = null;
    switch (schema.typeName) {
      case 'entity':
        decoded = Entity.fromBinary(new Uint8Array(raw)) as ProtoTypes;
        break;
      default:
        break;
    }

    await handler(decoded);
  }
}
