/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface EntitiesContentIntegration {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string
  /** @example true */
  enabled: boolean
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  id: string
  /** @example "FAQ" */
  name: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  project_id: string
  /** @example "Configurable floating button for your website" */
  summary: string
  /** @example "SuperButton is the best app to create configurable floating buttons on your website." */
  text: string
  /** @example "What is SuperButton?" */
  title: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  user_id: string
}

export interface EntitiesProject {
  /** @example "#283593" */
  color: string
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string
  /** @example "Need some help?" */
  greeting: string
  /** @example 0 */
  greeting_timeout_seconds: number
  /** @example "https://cdn.superbutton.app/chat-icon.svg" */
  icon: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  id: string
  /** @example "Joe's Store" */
  name: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string
  /** @example "https://example.com" */
  url: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  user_id: string
}

export interface EntitiesProjectIntegration {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  id: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  integration_id: string
  name: string
  /** @example 1 */
  position: number
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  project_id: string
  /** @example "whatsapp" */
  type: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  user_id: string
}

export interface EntitiesUser {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string
  /** @example "name@email.com" */
  email: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  id: string
  /** @example "John Doe" */
  name: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string
}

export interface EntitiesWhatsappIntegration {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string
  /** @example true */
  enabled: boolean
  /** @example "https://cdn.superbutton.app/whatsapp-icon.svg" */
  icon: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  id: string
  /** @example "FAQ" */
  name: string
  /** @example "+18005550199" */
  phone_number: string
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  project_id: string
  /** @example "Contact us on WhatsApp" */
  text: string
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  user_id: string
}

export interface RequestsCloudEvent {
  data: any
  datacontenttype: string
  id: string
  source: string
  specversion: string
  time: string
  type: string
}

export interface RequestsContentIntegrationCreateRequest {
  name: string
  summary: string
  text: string
  title: string
}

export interface RequestsContentIntegrationUpdateRequest {
  name: string
  summary: string
  text: string
  title: string
}

export interface RequestsProjectCreateRequest {
  name: string
  website: string
}

export interface RequestsProjectUpdateRequest {
  color: string
  greeting: string
  greeting_timeout: number
  icon: string
  name: string
  website: string
}

export interface RequestsWhatsappIntegrationCreateRequest {
  name: string
  phone_number: string
  text: string
}

export interface RequestsWhatsappIntegrationUpdateRequest {
  name: string
  phone_number: string
  text: string
}

export interface ResponsesBadRequest {
  /** @example "The request body is not a valid JSON string" */
  data: string
  /** @example "The request isn't properly formed" */
  message: string
  /** @example "error" */
  status: string
}

export interface ResponsesInternalServerError {
  /** @example "We ran into an internal error while handling the request." */
  message: string
  /** @example "error" */
  status: string
}

export interface ResponsesNoContent {
  /** @example "phone deleted successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesNotFound {
  /** @example "cannot find message with ID [32343a19-da5e-4b1b-a767-3298a73703ca]" */
  message: string
  /** @example "error" */
  status: string
}

export interface ResponsesOkArrayEntitiesProject {
  data: EntitiesProject[]
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesOkArrayEntitiesProjectIntegration {
  data: EntitiesProjectIntegration[]
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesOkEntitiesContentIntegration {
  data: EntitiesContentIntegration
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesOkEntitiesProject {
  data: EntitiesProject
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesOkEntitiesUser {
  data: EntitiesUser
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesOkEntitiesWhatsappIntegration {
  data: EntitiesWhatsappIntegration
  /** @example "Request handled successfully" */
  message: string
  /** @example "success" */
  status: string
}

export interface ResponsesUnauthorized {
  /** @example "Make sure your API key is set in the [X-API-Key] header in the request" */
  data: string
  /** @example "You are not authorized to carry out this request." */
  message: string
  /** @example "error" */
  status: string
}

export interface ResponsesUnprocessableEntity {
  data: Record<string, string[]>
  /** @example "validation errors while sending message" */
  message: string
  /** @example "error" */
  status: string
}
