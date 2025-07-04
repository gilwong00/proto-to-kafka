import { Module } from '@nestjs/common';
import { KafkaService } from './kafka.service';
import { MessageRouterService } from './message-router.service';

@Module({
  providers: [KafkaService, MessageRouterService],
})
export class KafkaModule {}
