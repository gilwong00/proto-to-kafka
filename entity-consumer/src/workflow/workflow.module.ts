import { Module } from '@nestjs/common';
import { EntityWorkflowService } from './entity.workflow.service';

@Module({
  providers: [EntityWorkflowService],
})
export class WorkflowModule {}
