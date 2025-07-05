/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-call */
import { Controller, Logger } from '@nestjs/common';
import { Ctx, KafkaContext, MessagePattern } from '@nestjs/microservices';
import { Entity } from 'gen/entity';
import { EntityService } from './entity.service';
import { EntityEventType } from './events';

@Controller('entity')
export class EntityController {
  private readonly logger = new Logger(EntityController.name);

  constructor(private readonly entityService: EntityService) {}

  @MessagePattern('entity-topic')
  handleEntityEvent(@Ctx() context: KafkaContext) {
    const topic = context.getTopic();
    const partition = context.getPartition();
    const message = context.getMessage();
    const offset = message.offset;
    const key = message.key?.toString();
    const headers = message.headers;

    this.logger.log(
      `Received message on ${topic} [${partition}] offset ${offset} key ${key}`,
    );

    const buffer = Buffer.isBuffer(message.value)
      ? message.value
      : Buffer.from(message.value!);

    const decodedEntity = Entity.fromBinary(buffer) as Entity;
    const event = (headers?.['eventType'] as EntityEventType) ?? '';
    this.entityService.processEntityEvent(event, decodedEntity);
  }
}
