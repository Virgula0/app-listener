# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 471 files · ~1,691,740 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6176 nodes · 26263 edges · 223 communities (186 shown, 37 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2479 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6bef1458`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- selection.go
- policyEditor
- mutedText
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- os.FileMode
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- TrustGuard
- guard.bpf.c
- NetEventType
- watchSet
- .FilterExistingWhitelist
- User
- UI
- Serve
- IntegrationSuite
- procstats.go
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- hashFileShared
- GuardEvent
- github.com/charmbracelet/bubbles/textarea.Model
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- Client
- watchGroup
- tabs.go
- Open
- RunWith
- icon
- GuardInodeKey
- filevault_test.go
- update_test.go
- newTestTarget
- exe_supersede.h
- Network
- SanitizeText
- bpfstats/main.go
- gui/session.go
- runMonitor
- filevault.go
- fakeGuardRepo
- fcntl
- runBackend
- Guard
- FileTarget
- runNetworkMonitor
- lockAndDeprovision
- updatePage
- RefreshSection
- time.Time
- process_vm_readv.c
- Monitor
- journal
- planUpdaters
- WithBottomBar
- GlobReservations
- CandidateDir
- codeedit.go
- DaemonEvent
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- testLimiter
- changelog.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- catalogRefresher
- fileedit_symlink_test.go
- update/update.go
- Well-known bypass classes
- BinaryEntry
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- installPage
- prepareDaemonStart
- mountTable
- TestMonitorUseCaseLifecycle
- Config
- IntegrationSuite
- IntegrationSuite
- net_tester/main.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- expandPlaceholders
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- SampleContent
- make test-integration (rootful Docker suite)
- headless.go
- supersede_probe.c
- go_pkg_fyne_io_fyne_v2
- trace-app-libs.sh
- IntegrationSuite
- sanitizeTerminalText
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- time.Duration
- IntegrationSuite
- runOfflineEdit
- hashCache
- runDaemon
- Execute
- SelectOrphans
- downloadFile
- NewEventFanout
- ConfinedEntry
- EventType
- layout.go
- copyRegularAt
- sync.Mutex
- Installation
- backend_test.go
- graphify-refresh.sh
- Vault
- inodeChain
- configedit_test.go
- pinstate.go
- feed.go
- shell
- runView
- makeReloadHandler
- runGenKey
- Finder
- FindDaemonProcesses
- Limitations
- Troubleshooting
- IntegrationSuite
- resolveFilesystemPrereqs
- Daemon configuration
- NewDumpHook
- .populateWatchInodes
- systemd.go
- Compatibility
- How it works
- fakeDaemon
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- fakeVault
- writeEvent
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- validateFlags
- pathCache
- run
- writeWithin
- NetGuardEvent
- parseGuardEvent
- Desktop GUI
- IntegrationSuite
- Ledger
- .eventTypeAllowed
- WithSuspectsFrom
- DiscoverInfraBinaries
- CleanupStalePins
- buildOneGuard
- Vault
- checkRunningBinaryMatchesInstalled
- golang.org/x/sys/unix.Timespec
- bleedLayout
- TestMain
- MulticallError
- partitionEncryptionRoots

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `GuardInodeKey` - 99 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (223 total, 37 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (16): tempJournalRow, invalidPassword, launchRule, trustEvent, btrfsIoctlFsInfoArgs, infraContainerPaths(), TestIntegrationSuite(), daemonStatsMsg (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (143): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), openTemp(), TestAuthSucceededResetsBudget() (+135 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.15
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (22): systemPatterns(), systemRule(), addErrKind, ReplacementCheck, Line(), BinaryRule, classifyAddErr(), addRecorder() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 10 - "selection.go"
Cohesion: 0.11
Nodes (17): pos, diffModel, Highlighter, cellToRune(), clipboardText(), cursorPos(), Selection, offset() (+9 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 12 - "mutedText"
Cohesion: 0.06
Nodes (23): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, Event, scrollable() (+15 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (10): editorModel, clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel (+2 more)

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "os.FileMode"
Cohesion: 0.11
Nodes (24): entryOf(), dupFD(), installConfig(), installFile(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+16 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (10): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), chownUser, fileEditMode, fileEditModel (+2 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (10): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth_test.go"
Cohesion: 0.14
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "TrustGuard"
Cohesion: 0.07
Nodes (15): tempGrant, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), loadPinnedExeMaps(), stripPinnedAllow() (+7 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "NetEventType"
Cohesion: 0.08
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+10 more)

### Community 27 - "watchSet"
Cohesion: 0.18
Nodes (11): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+3 more)

### Community 28 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 29 - "User"
Cohesion: 0.10
Nodes (27): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup() (+19 more)

### Community 30 - "UI"
Cohesion: 0.09
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), MultiSelectKeymap() (+11 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 33 - "procstats.go"
Cohesion: 0.09
Nodes (23): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+15 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+23 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.14
Nodes (19): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+11 more)

### Community 37 - "backend"
Cohesion: 0.09
Nodes (16): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), File, User, backend, backendOp (+8 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (26): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+18 more)

### Community 40 - "os.File"
Cohesion: 0.03
Nodes (91): setConfinedHomes(), VettedInode, BinaryStat, Multicall, MulticallKind, placedWalk, Start(), verifiedBinary() (+83 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (5): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "hashFileShared"
Cohesion: 0.23
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 45 - "GuardEvent"
Cohesion: 0.06
Nodes (25): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapGuard (+17 more)

### Community 46 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), newEditView(), relOf()

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+8 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (58): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), askEncryption(), libraryBlocksFromCandidates(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse() (+50 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (14): Client, launch(), TestReaderCapsLines(), Hello, Reader, Writer, NewReader(), hashFile() (+6 more)

### Community 53 - "watchGroup"
Cohesion: 0.08
Nodes (43): grantNewBinaries(), watchGroup, Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource(), applyAllowLib(), applyConfigLine(), applyDirective() (+35 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (15): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+7 more)

### Community 55 - "Open"
Cohesion: 0.14
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 56 - "RunWith"
Cohesion: 0.15
Nodes (22): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups() (+14 more)

### Community 57 - "icon"
Cohesion: 0.06
Nodes (30): binaryList, binaryRow, clickPanel, counterStrip, filesPage, handButton, homePage, iconTile (+22 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (54): engine, refusedReplacements, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), appendInodeRow() (+46 more)

### Community 59 - "filevault_test.go"
Cohesion: 0.17
Nodes (27): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestEnsureRecoverySidecarPlaceholder() (+19 more)

### Community 60 - "update_test.go"
Cohesion: 0.08
Nodes (21): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchGUIWithoutAsset(), TestFetchReleaseByTag(), TestFetchReleases() (+13 more)

### Community 61 - "newTestTarget"
Cohesion: 0.24
Nodes (16): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "Network"
Cohesion: 0.12
Nodes (9): netGuardPill(), newExeResolver(), guardRow(), clonePolicy(), NetExtra, NetPolicy, Network, exeCached (+1 more)

### Community 64 - "SanitizeText"
Cohesion: 0.13
Nodes (13): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), startTrustGuard(), trustStartupError() (+5 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.06
Nodes (37): CheckBPFLSM(), runBPFCheck(), BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt() (+29 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): invokingIDs(), ancestorEnv(), cleanupSessionEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "runMonitor"
Cohesion: 0.09
Nodes (21): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets(), validatePair() (+13 more)

### Community 68 - "filevault.go"
Cohesion: 0.19
Nodes (16): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+8 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "runBackend"
Cohesion: 0.08
Nodes (30): auditAfterEdit(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply(), editAgain() (+22 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (23): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), binaryVerifyState, deferredBinary, exeRow, rootHandle, tempMaskOp (+15 more)

### Community 73 - "FileTarget"
Cohesion: 0.12
Nodes (13): Target, guardPill(), RealFileEngines(), FileEngines, Files, FileSpec, FileTarget, Policy (+5 more)

### Community 74 - "runNetworkMonitor"
Cohesion: 0.12
Nodes (10): CheckEBPF(), ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository (+2 more)

### Community 75 - "lockAndDeprovision"
Cohesion: 0.11
Nodes (13): deprovisionKind, migrationTarget, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+5 more)

### Community 76 - "updatePage"
Cohesion: 0.22
Nodes (6): installedBinary, latestRelease, updatePage, readInstalled(), boolIndex(), isPreRelease()

### Community 77 - "RefreshSection"
Cohesion: 0.08
Nodes (32): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+24 more)

### Community 78 - "time.Time"
Cohesion: 0.15
Nodes (11): LockedError, newWatchSet(), fakeInotify, inodeBatchFn, sweepID, sweepMemo, fakeBatches(), TestBatchInodeRows() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 81 - "journal"
Cohesion: 0.14
Nodes (5): fakeGuard, fakeMonitor, fakeNet, fakeNetMonitor, journal

### Community 82 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 83 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 84 - "GlobReservations"
Cohesion: 0.15
Nodes (11): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+3 more)

### Community 85 - "CandidateDir"
Cohesion: 0.10
Nodes (31): selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories(), groupCandidates(), pickDirectories(), pickFromCandidates(), pickUsers(), TestGroupCandidates() (+23 more)

### Community 86 - "codeedit.go"
Cohesion: 0.09
Nodes (44): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+36 more)

### Community 87 - "DaemonEvent"
Cohesion: 0.22
Nodes (13): ServeConfig, newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats() (+5 more)

### Community 88 - "liveSession"
Cohesion: 0.22
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "testLimiter"
Cohesion: 0.18
Nodes (16): drainEphemeral(), runHeadless(), admitEvent(), foldable(), newGateLogLimiter(), targetProgram(), gateEvent(), TestGateLogLimiterFlushesOnShutdown() (+8 more)

### Community 92 - "changelog.go"
Cohesion: 0.16
Nodes (23): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+15 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (21): NewMountVouch(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv() (+13 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "catalogRefresher"
Cohesion: 0.16
Nodes (7): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 100 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 101 - "update/update.go"
Cohesion: 0.10
Nodes (44): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit() (+36 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "BinaryEntry"
Cohesion: 0.14
Nodes (14): keptEntries, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath(), diffBinaries(), snapEntry(), BinariesSummary() (+6 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.11
Nodes (14): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, UninstallHooks (+6 more)

### Community 110 - "prepareDaemonStart"
Cohesion: 0.14
Nodes (17): configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), noteReloadedSettings(), prepareDaemonStart(), resolveConfigPath(), runLockdown(), runOneShotMode() (+9 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.16
Nodes (14): NewGuardUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+6 more)

### Community 113 - "Config"
Cohesion: 0.06
Nodes (49): newCatalogRefresher(), startCatalogRefresh(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure() (+41 more)

### Community 116 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "expandPlaceholders"
Cohesion: 0.14
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "SampleContent"
Cohesion: 0.12
Nodes (20): TestDaemonUnitWithMetadataOutput(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent(), SampleFiles(), TestDaemonUnitOrdersUserCodeAfterReadiness(), TestNoLegacyCatalogRefreshSamples() (+12 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "go_pkg_fyne_io_fyne_v2"
Cohesion: 0.05
Nodes (33): buttonCursor, fakeLink, multiAnswer, pageDef, textAnswer, uninstallCard, handCursors(), handCursorsOn() (+25 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.14
Nodes (14): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant(), newExeTap() (+6 more)

### Community 143 - "time.Duration"
Cohesion: 0.14
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 145 - "runOfflineEdit"
Cohesion: 0.12
Nodes (16): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+8 more)

### Community 146 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 147 - "runDaemon"
Cohesion: 0.18
Nodes (12): awaitStartupOrSignal(), catchLifecycleSignals(), markReady(), notifySystemdReady(), runDaemon(), startGuardedDaemonAbortable(), startPprof(), writePidFile() (+4 more)

### Community 148 - "Execute"
Cohesion: 0.20
Nodes (9): releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand() (+1 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 150 - "downloadFile"
Cohesion: 0.33
Nodes (4): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 152 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (18): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+10 more)

### Community 154 - "EventType"
Cohesion: 0.09
Nodes (15): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+7 more)

### Community 155 - "layout.go"
Cohesion: 0.22
Nodes (10): adaptiveGrid, checkGrid, gapVBox, TestAdaptiveGridDropsColumns(), cellMin(), newAdaptiveGrid(), newCheckGrid(), newGapVBox() (+2 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "sync.Mutex"
Cohesion: 0.22
Nodes (6): newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.16
Nodes (26): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+18 more)

### Community 161 - "Vault"
Cohesion: 0.07
Nodes (35): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+27 more)

### Community 162 - "inodeChain"
Cohesion: 0.18
Nodes (18): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent(), statFD() (+10 more)

### Community 163 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 164 - "pinstate.go"
Cohesion: 0.13
Nodes (27): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+19 more)

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (37): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, networkPage (+29 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (34): pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet, sysStatus (+26 more)

### Community 167 - "runView"
Cohesion: 0.06
Nodes (15): logGrid, logHook, logLine, runHost, runView, stageItem, stageState, logRow() (+7 more)

### Community 168 - "makeReloadHandler"
Cohesion: 0.29
Nodes (4): newControlManager(), startControlManager(), makeReloadHandler(), controlManager

### Community 169 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 170 - "Finder"
Cohesion: 0.13
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 171 - "FindDaemonProcesses"
Cohesion: 0.17
Nodes (12): anyDaemonRunning(), probeSystem(), Check(), collectDiagnostics(), readSysctl(), requireKernel51(), DaemonRunning(), FindDaemonProcesses() (+4 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "resolveFilesystemPrereqs"
Cohesion: 0.19
Nodes (9): encryptDirectories(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54() (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 179 - "systemd.go"
Cohesion: 0.06
Nodes (58): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), cleanupBackups(), deploy() (+50 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 195 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 196 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "validateFlags"
Cohesion: 0.33
Nodes (5): TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), TestUnchangedPerShared()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "run"
Cohesion: 0.15
Nodes (17): confirmKeyDeletion(), run(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), TestRunWithRefusesConcurrentRun(), InstallCert (+9 more)

### Community 202 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 204 - "NetGuardEvent"
Cohesion: 0.08
Nodes (22): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNetGuard, eventsetSummary() (+14 more)

### Community 205 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.11
Nodes (12): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, applySchema() (+4 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 213 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 215 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 216 - "buildOneGuard"
Cohesion: 0.08
Nodes (37): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+29 more)

### Community 218 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 567 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **37 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `RunWith`, `backend`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09721499376671791 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019103531823688378 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `RunWith` to `daemon/daemon.go`, `edit-protected`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._