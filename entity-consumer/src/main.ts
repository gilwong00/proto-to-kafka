import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { ConfigService } from '@nestjs/config';
import { logLevel } from '@nestjs/microservices/external/kafka.interface';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const configService = app.get(ConfigService);
  const port = configService.get<string>('PORT') ?? 3333;

  // Create Kafka microservice
  const brokers = configService
    .get<string>('KAFKA_BROKERS', 'localhost:9092')
    .split(',');

  app.connectMicroservice<MicroserviceOptions>({
    transport: Transport.KAFKA,
    options: {
      client: {
        brokers,
        // TODO: this should be configurable in a production system
        logLevel: logLevel.INFO,
      },
    },
  });

  await app.startAllMicroservices();
  await app.listen(port);
}
bootstrap();
