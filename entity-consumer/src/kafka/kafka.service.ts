import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Kafka } from 'kafkajs';
import * as protobuf from 'protobufjs';
import * as path from 'path';

@Injectable()
export class KafkaService implements OnModuleInit {
  private readonly logger = new Logger(KafkaService.name);
  private kafkaClient: Kafka;

  constructor(private readonly configService: ConfigService) {
    const kafkaEndpoint: string =
      this.configService.get('KAKFA_ENDPOINT') ?? '';

    this.kafkaClient = new Kafka({ brokers: [kafkaEndpoint] });
  }

  async onModuleInit() {}
}
