export const EntityEvents = {
  ENTITY_CREATED: 'entity-created',
} as const;

export type EntityEventType = (typeof EntityEvents)[keyof typeof EntityEvents];
