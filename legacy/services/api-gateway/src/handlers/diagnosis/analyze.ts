import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { BedrockRuntimeClient, InvokeModelCommand } from '@aws-sdk/client-bedrock-runtime';
import { handleError, successResponse } from '../../middleware/error-handler';
import { z } from 'zod';

const bedrockClient = new BedrockRuntimeClient({ region: process.env.AWS_REGION ?? 'ap-southeast-1' });

const diagnosisSchema = z.object({
  symptoms: z.array(z.string()).min(1).max(10),
  vehicleMake: z.string().optional(),
  vehicleModel: z.string().optional(),
  vehicleYear: z.number().optional(),
  obdCodes: z.array(z.string()).optional(),
  photos: z.array(z.string()).optional(),
});

const SYSTEM_PROMPT = `You are a Filipino automotive diagnostic assistant for TowCommand PH.
Analyze the reported vehicle symptoms and provide:
1. Likely issue diagnosis (ranked by probability)
2. Recommended service type: FLATBED_TOW, WHEEL_LIFT, JUMPSTART, TIRE_CHANGE, FUEL_DELIVERY, LOCKOUT, or ACCIDENT_RECOVERY
3. Urgency level: low, medium, high, critical
4. Whether the vehicle is safe to drive
5. Estimated repair cost range in PHP
Keep responses concise and practical. Use simple English that Filipino users can understand.`;

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const body = diagnosisSchema.parse(JSON.parse(event.body ?? '{}'));

    const userMessage = [
      `Symptoms: ${body.symptoms.join(', ')}`,
      body.vehicleMake ? `Vehicle: ${body.vehicleMake} ${body.vehicleModel ?? ''} ${body.vehicleYear ?? ''}` : '',
      body.obdCodes?.length ? `OBD Codes: ${body.obdCodes.join(', ')}` : '',
    ].filter(Boolean).join('\n');

    const response = await bedrockClient.send(new InvokeModelCommand({
      modelId: 'anthropic.claude-3-5-sonnet-20241022-v2:0',
      contentType: 'application/json',
      accept: 'application/json',
      body: JSON.stringify({
        anthropic_version: 'bedrock-2023-05-31',
        max_tokens: 1024,
        system: SYSTEM_PROMPT,
        messages: [{ role: 'user', content: userMessage }],
      }),
    }));

    const result = JSON.parse(new TextDecoder().decode(response.body));
    const analysisText = result.content?.[0]?.text ?? 'Unable to analyze symptoms';

    return successResponse({
      analysis: analysisText,
      symptoms: body.symptoms,
      obdCodes: body.obdCodes,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    return handleError(error);
  }
}
