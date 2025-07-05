import { Module } from '@nestjs/common';
import { EntityService } from './entity.service';
import { EntityController } from './entity.controller';

@Module({
  providers: [EntityService],
  controllers: [EntityController]
})
export class EntityModule {}
