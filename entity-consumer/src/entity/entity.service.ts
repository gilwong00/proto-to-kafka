import { Injectable, Logger } from '@nestjs/common';
import { Entity } from 'gen/entity';
import { EntityEvents, EntityEventType } from './events';

@Injectable()
export class EntityService {
  private readonly logger = new Logger(EntityService.name);

  processEntityEvent(event: EntityEventType, entity: Entity) {
    switch (event) {
      case EntityEvents.ENTITY_CREATED:
        this.logger.log('New entity created', JSON.stringify(entity));
        break;
      default:
        break;
    }
  }
}
