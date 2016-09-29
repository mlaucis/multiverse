package com.tapglue.signals;

import com.google.api.services.bigquery.model.TableFieldSchema;
import com.google.api.services.bigquery.model.TableRow;
import com.google.api.services.bigquery.model.TableSchema;
import com.google.cloud.dataflow.sdk.Pipeline;
import com.google.cloud.dataflow.sdk.coders.Coder;
import com.google.cloud.dataflow.sdk.coders.protobuf.ProtoCoder;
import com.google.cloud.dataflow.sdk.io.BigQueryIO;
import com.google.cloud.dataflow.sdk.io.PubsubIO;
import com.google.cloud.dataflow.sdk.options.DataflowPipelineOptions;
import com.google.cloud.dataflow.sdk.options.Default;
import com.google.cloud.dataflow.sdk.options.DefaultValueFactory;
import com.google.cloud.dataflow.sdk.options.Description;
import com.google.cloud.dataflow.sdk.options.PipelineOptions;
import com.google.cloud.dataflow.sdk.options.PipelineOptionsFactory;
import com.google.cloud.dataflow.sdk.runners.DataflowPipelineRunner;
import com.google.cloud.dataflow.sdk.transforms.DoFn;
import com.google.cloud.dataflow.sdk.transforms.ParDo;

import com.tapglue.signals.data.Protobufs;
import com.tapglue.signals.data.Protobufs.Signal;

import java.util.ArrayList;

public class SignalsPersist {
	static class BigQueryConverter extends DoFn<Signal, TableRow> {
		@Override
		public void processElement(ProcessContext c) {
			TableRow row = new TableRow()
					.set("id", c.element().getId())
                    .set("payload", c.element().toByteArray());

            c.output(row);
		}

		static TableSchema getSchema() {
			return new TableSchema().setFields(new ArrayList<TableFieldSchema>() {
				{
					add(new TableFieldSchema().setName("id").setType("INTEGER").setMode("REQUIRED"));
					add(new TableFieldSchema().setName("payload").setType("BYTES").setMode("REQUIRED"));
				}
			});
		}
	}

	public static interface SignalsPersistOptions extends DataflowPipelineOptions {
		@Description("")
		@Default.InstanceFactory(TopicFactory.class)
		String getPubsubTopic();
		void setPubsubTopic(String topic);

		static class TopicFactory implements DefaultValueFactory<String> {
			@Override
			public String create(PipelineOptions options) {
				DataflowPipelineOptions dataflowPipelineOptions =
					options.as(DataflowPipelineOptions.class);

				return "projects/tapglue-signals/topics/test-events";
			}
		}
	}

	public static void main(String[] args) {
		SignalsPersistOptions options = PipelineOptionsFactory.fromArgs(args)
			.withValidation()
			.as(SignalsPersistOptions.class);

		options.setRunner(DataflowPipelineRunner.class);
		options.setStreaming(true);

		String tableSpec = new StringBuilder()
	        .append(options.getProject()).append(":")
	        .append("signals_persist").append(".")
	        .append("all")
	        .toString();
		Coder<Signal> coder = ProtoCoder.of(Signal.class).withExtensionsFrom(Protobufs.class);
		Pipeline pipeline = Pipeline.create(options);

		pipeline
			.apply(PubsubIO.Read.named("ConsumePubSub").topic(options.getPubsubTopic()).withCoder(coder))
			.apply(ParDo.of(new BigQueryConverter()))
			.apply(BigQueryIO.Write
				.named("Persist")
				.to(tableSpec)
				.withSchema(BigQueryConverter.getSchema())
				.withWriteDisposition(BigQueryIO.Write.WriteDisposition.WRITE_APPEND)
				.withCreateDisposition(BigQueryIO.Write.CreateDisposition.CREATE_IF_NEEDED));

		pipeline.run();
	}
}